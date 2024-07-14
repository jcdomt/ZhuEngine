package main

import (
	"ZhuEngine/config"
	"ZhuEngine/router"
	"net/http"
	"strconv"

	"log"
)

func main() {
	// 读取配置
	conf := config.GetConfig()

	router.LoadSitesRouter(conf)

	log.Default().Print("主程序启动中")
	http.Handle("/", &router.Pxy{})
	port_string := strconv.Itoa(conf.ZhuEngine.Port)
	http.ListenAndServe("0.0.0.0:"+port_string, nil)
}
