package config

type Config struct {
	ZhuEngine *ZeConfig
	HTTPS     *HttpsConfig
	Web       *WebConfig
	Cgi       map[string]*CgiConfig
}

func GetConfig() *Config {
	conf := new(Config)
	// 获取 ZhuEngine 主要配置

	conf.ZhuEngine = GetZeConfig()
	conf.HTTPS = GetHttpsConfig()
	conf.Web = GetWebConfig()
	conf.Cgi = make(map[string]*CgiConfig)
	conf.Cgi = getCgiConfig()
	return conf
}
