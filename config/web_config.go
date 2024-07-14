package config

import "gopkg.in/ini.v1"

type SiteConfig struct {
	Type    string `ini:"type"`    // 类型：domain 子域名 url 子路劲
	Url     string `ini:"url"`     // 子域名或子路径
	Server  string `ini:"server"`  // 对应服务主项目
	Autorun bool   `ini:"autorun"` // 是否自动拉起服务
	Exec    string `ini:"exec"`    // 自动拉起服务的指令
}
type WebConfig struct {
	Sites map[string]*SiteConfig
}

func GetWebConfig() *WebConfig {
	wc := new(WebConfig)
	wc.Sites = make(map[string]*SiteConfig)
	inicfg, err := ini.Load("./conf/web.ini")
	if err != nil {
		panic(err)
	}
	sections := inicfg.Sections()
	for _, v := range sections {
		if v.Name() == "DEFAULT" {
			continue
		}
		m := new(SiteConfig)
		v.MapTo(m)
		wc.Sites[v.Name()] = m
	}
	return wc
}
