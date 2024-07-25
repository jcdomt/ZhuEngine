package main

import (
	"ZhuEngine/config"
	"ZhuEngine/router"
	"ZhuEngine/site"
	"net/http"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	// 读取配置
	conf := config.GetConfig()

	sites := site.LoadSitesRouter(conf)
	site.SiteAutoRun(sites)

	log.Info("启动主程序")
	handler := &router.Pxy{}
	http.Handle("/", handler)
	port_string := strconv.Itoa(conf.ZhuEngine.Port)
	https_port_string := strconv.Itoa(conf.HTTPS.Port)
	if !conf.HTTPS.Enable {
		http.ListenAndServe("0.0.0.0:"+port_string, handler)
	} else {
		go http.ListenAndServe("0.0.0.0:"+port_string, handler)
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+https_port_string, conf.HTTPS.Crt, conf.HTTPS.Key, handler))
	}

}
