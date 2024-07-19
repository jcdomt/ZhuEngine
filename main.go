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
	http.Handle("/", &router.Pxy{})
	port_string := strconv.Itoa(conf.ZhuEngine.Port)
	http.ListenAndServe("0.0.0.0:"+port_string, nil)
}
