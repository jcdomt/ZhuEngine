// 核心代理功能
package router

import (
	"io"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := req.Host + req.URL.Path
	site := getRequestSite(url)

	if site == nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	r := site.SendHttp(rw, req)
	if r != nil {
		io.Copy(rw, r.Body)
	} else {
		rw.WriteHeader(http.StatusBadGateway)
	}
}

func getRequestSite(url string) *Site {
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "http://", "", 1)
	// 先判断子域名项目
	sub_domain := strings.Split(url, ".")
	value, exist := Sites_SubDomain[sub_domain[0]]
	if exist {
		return value
	} else {
		sub_patter := strings.Split(url, "/")
		value, exist := Sites_SubPatter[sub_patter[1]]
		if exist {
			return value
		} else {
			return nil
		}
	}
}
