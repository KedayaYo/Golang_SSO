package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// 定义全局配置变量
var cfg App

// 设置配置文件
func Setup() {
	// 定义命令行参数，用于指定配置文件路径
	path := flag.String("config", "/etc/oauth2nsso/config.yaml", "config.yaml的绝对路径")
	// 解析命令行参数
	flag.Parse()

	// 读取配置文件内容
	content, err := os.ReadFile(*path)
	if err != nil {
		// 读取文件出错，打印错误信息并退出程序
		log.Fatalf("error: %v", err)
	}

	// 解析YAML格式的配置文件内容到全局配置变量
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		// 解析YAML内容出错，打印错误信息并退出程序
		log.Fatalf("error: %v", err)
	}
}
