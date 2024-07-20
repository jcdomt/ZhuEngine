// 核心代理功能
package router

import (
	"ZhuEngine/site"
	"net/http"
	"path/filepath"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := req.Host + req.URL.Path
	s := getRequestSite(url)

	if s == nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "text/html; charset=utf-8")
	}
	if s.CgiEnable {
		// 启用了 CGI 功能
		filename := strings.TrimPrefix(req.URL.Path, "/")
		if filename == "" {
			// 使用默认 CGI 文件名
			filename = s.CgiDefaultFilename
		}
		scriptPath := filepath.Join(s.Config.Server, filename)
		handler := s.GenerateCgiHandler(scriptPath)
		handler.ServeHTTP(rw, req)
	} else {
		s.SendHttp(rw, req)
	}

	// r := s.SendHttp(rw, req)
	// if r != nil {
	// 	io.Copy(rw, r.Body)
	// } else {
	// 	rw.WriteHeader(http.StatusBadGateway)
	// }
}

func getRequestSite(url string) *site.Site {
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "http://", "", 1)
	// 先判断子域名项目
	sub_domain := strings.Split(url, ".")
	value, exist := site.Sites_SubDomain[sub_domain[0]]
	if exist {
		return value
	} else {
		sub_patter := strings.Split(url, "/")
		value, exist := site.Sites_SubPatter[sub_patter[1]]
		if exist {
			return value
		} else {
			return nil
		}
	}
}
