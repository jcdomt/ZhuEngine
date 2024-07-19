package config

import (
	"gopkg.in/ini.v1"
)

type CgiConfig struct {
	CGI     string `ini:"cgi"` // 相关 cgi 的目录
	Default string `ini:"default"`
}

func getCgiConfig() map[string]*CgiConfig {
	ret := make(map[string]*CgiConfig)

	inicfg, err := ini.Load("./conf/cgi.ini")
	if err != nil {
		panic(err)
	}
	sections := inicfg.Sections()
	for _, v := range sections {
		if v.Name() == "DEFAULT" {
			continue
		}
		m := new(CgiConfig)
		v.MapTo(m)
		ret[v.Name()] = m
	}

	return ret
}
