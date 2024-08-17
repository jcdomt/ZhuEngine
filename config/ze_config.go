package config

type ZeConfig struct {
	Port int    `ini:"port"` // 核心服务端口
	Host string `ini:"host"` // 网站服务URL
}

func GetZeConfig() *ZeConfig {
	inicfg, _ := parseIniConfigSyntax("./conf/ze.ini")

	ze_config := new(ZeConfig)
	err := inicfg.Section("zhu_engine").MapTo(ze_config)
	if err != nil {
		panic(err)
	}
	return ze_config
}

type HttpsConfig struct {
	Enable bool   `ini:"enable"`
	Crt    string `ini:"crt"`
	Key    string `ini:"key"`
	Port   int    `ini:"port"`
	Force  bool   `ini:"force"`
}

func GetHttpsConfig() *HttpsConfig {
	inicfg, _ := parseIniConfigSyntax("./conf/ze.ini")

	https_config := new(HttpsConfig)
	err := inicfg.Section("https").MapTo(https_config)
	if err != nil {
		panic(err)
	}

	return https_config
}
