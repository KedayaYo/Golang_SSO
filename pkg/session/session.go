package session

import (
	"encoding/gob"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/llaoj/oauth2nsso/config"
)

// 声明全局的会话存储变量
var store *sessions.CookieStore

// 设置会话存储配置
func Setup() {
	// 注册类型，以便在会话中存储该类型
	gob.Register(url.Values{})

	// 初始化Cookie会话存储
	store = sessions.NewCookieStore([]byte(config.Get().Session.SecretKey))
	store.Options = &sessions.Options{
		Path:     "/",                         // 设置Cookie路径为根路径
		MaxAge:   config.Get().Session.MaxAge, // 会话有效期，单位为秒
		HttpOnly: true,                        // 设置HttpOnly属性，提高安全性
	}
}

// 获取会话中的值
func Get(r *http.Request, name string) (val interface{}, err error) {
	// 获取会话
	session, err := store.Get(r, config.Get().Session.Name)
	if err != nil {
		return
	}

	// 获取指定键的值
	val = session.Values[name]

	return
}

// 设置会话中的值
func Set(w http.ResponseWriter, r *http.Request, name string, val interface{}) (err error) {
	// 获取会话
	session, err := store.Get(r, config.Get().Session.Name)
	if err != nil {
		return
	}

	// 设置指定键的值
	session.Values[name] = val
	// 保存会话
	err = session.Save(r, w)

	return
}

// 删除会话中的值
func Delete(w http.ResponseWriter, r *http.Request, name string) (err error) {
	// 获取会话
	session, err := store.Get(r, config.Get().Session.Name)
	if err != nil {
		return
	}

	// 删除指定键的值
	delete(session.Values, name)
	// 保存会话
	err = session.Save(r, w)

	return
}
