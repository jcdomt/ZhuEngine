// 负责站点的路由解析
package site

import (
	"ZhuEngine/config"
	"net/http/cgi"

	"github.com/go-kratos/kratos/v2/log"
)

type Site struct {
	Name      string // 站点名称
	Config    *config.SiteConfig
	Server    string
	URL       string
	SubDomain string
	SubPatter string
	Sub       string

	// CGI 功能
	CgiEnable          bool // 是否启用 CGI
	CgiHandler         *cgi.Handler
	CgiPath            string
	CgiDefaultFilename string

	// 负载均衡调度器功能
	ScheduleEnable bool      // 是否启用
	ScheduleType   string    // 调度手段
	Schedulor      Schedulor // 调度器实例
}

var Sites []*Site
var Sites_SubDomain map[string]*Site
var Sites_SubPatter map[string]*Site
var Site_RootDomain *Site

func init() {
	Sites = make([]*Site, 0)
	Sites_SubDomain = make(map[string]*Site)
	Sites_SubPatter = make(map[string]*Site)
}

func LoadSitesRouter(conf *config.Config) []*Site {
	sites, site1, site2 := LoadSites(conf)
	// 保存到全局变量
	Sites = sites
	for _, v := range site1 {
		Sites_SubDomain[v.SubDomain] = v
	}
	for _, v := range site2 {
		Sites_SubPatter[v.SubPatter] = v
	}
	return sites
}

func LoadSites(conf *config.Config) ([]*Site, []*Site, []*Site) {
	sites := make([]*Site, 0)
	sites_subdomain := make([]*Site, 0)
	sites_subpatter := make([]*Site, 0)
	for k, v := range conf.Web.Sites {
		site := new(Site)
		site.Name = k
		site.Config = v
		site.Server = v.Server
		switch site.Config.Type {
		case "domain":
			site.URL = site.Config.Url + "." + conf.ZhuEngine.Host
			site.SubDomain = site.Config.Url
			site.Sub = site.Config.Url
			sites_subdomain = append(sites_subdomain, site)
		case "patter":
			path := ""
			if site.Config.Url[0] != '/' {
				path = "/"
			}
			site.URL = conf.ZhuEngine.Host + path + site.Config.Url
			site.SubPatter = site.Config.Url
			site.Sub = site.Config.Url
			sites_subpatter = append(sites_subpatter, site)
		case "root_domain":
			site.URL = conf.ZhuEngine.Host
			Site_RootDomain = site
		}
		// CGI 功能判定
		if v.CGI != "" {
			// 存在 cgi 配置
			site.CgiEnable = true
			// 直接生成 CGI 处理器
			// handler := &cgi.Handler{
			// 	Path: conf.Cgi[v.CGI].CGI,     // 替换为实际的脚本路径
			// 	Dir:  v.Server + "/index.php", // CGI 脚本的 URL 路径前缀
			// }

			// handler := &cgi.Handler{
			// 	Path: conf.Cgi[v.CGI].CGI,
			// 	Dir:  filepath.Dir(v.Server),
			// 	Root: "/",
			// 	Env: []string{
			// 		"REDIRECT_STATUS=200",
			// 		"SCRIPT_FILENAME=" + v.Server, // 替换为实际的 PHP 脚本路径
			// 	},
			// }

			site.CgiPath = conf.Cgi[v.CGI].CGI
			site.CgiDefaultFilename = conf.Cgi[v.CGI].Default
			handler := &cgi.Handler{
				Path: site.CgiPath,
				Dir:  site.Config.Server,
				Root: "/",
			}
			site.CgiHandler = handler
		} else {
			site.CgiEnable = false
		}

		if v.Schedule != "" {
			site.ScheduleEnable = true
			site.ScheduleType = v.Schedule
			switch v.Schedule {
			case "round":
				schedulor := new(RoundRobinSchedulor)
				schedulor.Init(site)
				site.Schedulor = schedulor
			}
		}

		log.Info(site.Name, "启动成功：", site.Config.Server)
		sites = append(sites, site)
	}
	return sites, sites_subdomain, sites_subpatter
}
