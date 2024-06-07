package ldap

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/llaoj/oauth2nsso/config"
)

// 定义Session结构体
type Session struct {
	ldapCfg  config.LDAP // LDAP配置
	ldapConn *ldap.Conn  // LDAP连接
}

// 新建一个LDAP会话
func NewSession(ldapCfg config.LDAP) *Session {
	return &Session{
		ldapCfg: ldapCfg,
	}
}

// 格式化LDAP URL
func formatURL(ldapURL string) (string, error) {
	var protocol, hostport string
	_, err := url.Parse(ldapURL)
	if err != nil {
		return "", fmt.Errorf("解析LDAP主机错误: %s", err)
	}

	if strings.Contains(ldapURL, "://") {
		splitLdapURL := strings.Split(ldapURL, "://")
		protocol, hostport = splitLdapURL[0], splitLdapURL[1]
		if !((protocol == "ldap") || (protocol == "ldaps")) {
			return "", fmt.Errorf("未知的LDAP协议")
		}
	} else {
		hostport = ldapURL
		protocol = "ldap"
	}

	if strings.Contains(hostport, ":") {
		_, port, err := net.SplitHostPort(hostport)
		if err != nil {
			return "", fmt.Errorf("非法的LDAP URL, 错误: %v", err)
		}
		if port == "636" {
			protocol = "ldaps"
		}
	} else {
		switch protocol {
		case "ldap":
			hostport = hostport + ":389"
		case "ldaps":
			hostport = hostport + ":636"
		}
	}

	fLdapURL := protocol + "://" + hostport
	return fLdapURL, nil
}

// 打开LDAP会话
// 每次调用Open都应该调用Close
func (s *Session) Open() error {
	ldapURL, err := formatURL(s.ldapCfg.URL)
	if err != nil {
		return err
	}
	splitLdapURL := strings.Split(ldapURL, "://")

	protocol, hostport := splitLdapURL[0], splitLdapURL[1]
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return err
	}

	log.Println(ldapURL)

	switch protocol {
	case "ldap":
		l, err := ldap.Dial("tcp", hostport)
		if err != nil {
			return err
		}
		s.ldapConn = l
	case "ldaps":
		l, err := ldap.DialTLS("tcp", hostport, &tls.Config{ServerName: host, InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		s.ldapConn = l
	}

	return nil
}

// 关闭当前会话
func (s *Session) Close() {
	if s.ldapConn != nil {
		s.ldapConn.Close()
	}
}

// 用户认证
func UserAuthentication(username, password string) (string, error) {
	s := NewSession(config.Get().LDAP)
	if err := s.Open(); err != nil {
		return "", err
	}
	defer s.Close()

	// 首先用只读用户绑定
	if err := s.ldapConn.Bind(s.ldapCfg.SearchDN, s.ldapCfg.SearchPassword); err != nil {
		return "", err
	}

	// 根据用户名搜索用户
	searchRequest := ldap.NewSearchRequest(
		s.ldapCfg.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(s.ldapCfg.Filter, ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := s.ldapConn.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("用户不存在或者不唯一")
	}

	userdn := sr.Entries[0].DN

	// 绑定用户以验证其密码
	if err := s.ldapConn.Bind(userdn, password); err != nil {
		return "", err
	}

	// 重新绑定只读用户以进行进一步查询
	if err := s.ldapConn.Bind(s.ldapCfg.SearchDN, s.ldapCfg.SearchPassword); err != nil {
		return "", err
	}

	return username, nil
}
