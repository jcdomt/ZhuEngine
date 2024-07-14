package config

import (
	"gopkg.in/ini.v1"
)

type ZeConfig struct {
	Port int    `ini:"port"` // 核心服务端口
	Host string `ini:"host"` // 网站服务URL
}

func GetZeConfig() *ZeConfig {
	inicfg, _ := ini.Load("./conf/ze.ini")

	ze_config := new(ZeConfig)
	err := inicfg.Section("zhu_engine").MapTo(ze_config)
	if err != nil {
		panic(err)
	}
	return ze_config
}
