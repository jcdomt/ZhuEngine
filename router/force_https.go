package router

import (
	"net/http"
)

type ForceHttpsPxy struct{}

func (p *ForceHttpsPxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := "https://" + req.Host + req.URL.Path
	http.Redirect(rw, req, url, http.StatusTemporaryRedirect)

}
