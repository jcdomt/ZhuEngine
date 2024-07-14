package config

type Config struct {
	ZhuEngine *ZeConfig
	Web       *WebConfig
}

func GetConfig() *Config {
	conf := new(Config)
	// 获取 ZhuEngine 主要配置

	conf.ZhuEngine = GetZeConfig()
	conf.Web = GetWebConfig()
	return conf
}
