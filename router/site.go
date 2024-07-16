// 负责站点的路由解析
package router

import (
	"ZhuEngine/config"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Site struct {
	Name      string // 站点名称
	Config    *config.SiteConfig
	URL       string
	SubDomain string
	SubPatter string
	Sub       string
}

var Sites []*Site
var Sites_SubDomain map[string]*Site
var Sites_SubPatter map[string]*Site

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
		}
		sites = append(sites, site)
	}
	return sites, sites_subdomain, sites_subpatter
}

func (s *Site) SendHttp(rw http.ResponseWriter, req *http.Request) *http.Response {
	transport := http.DefaultTransport

	// step 1
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	// 正式的后台服务器地址
	//target := "http://" + s.Config.Server
	outReq.URL.Scheme = "http"
	outReq.URL.Host = s.Config.Server
	outReq.URL.Path = req.URL.Path
	outReq.URL.RawQuery = req.URL.RawQuery

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step 2
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		log.Default().Fatalln(err)
		rw.WriteHeader(http.StatusBadGateway)
		return nil
	}

	// step 3
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
	return res
}
