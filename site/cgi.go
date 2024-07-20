// 负责处理 http cgi 功能
// 其实除妖就是 php
package site

import (
	"net/http/cgi"
)

// 生成默认 CGI 处理器
func (s *Site) GenerateCgiHandler(scriptPath string) *cgi.Handler {
	defaultHandler := s.CgiHandler
	defaultHandler.Env = []string{
		"REDIRECT_STATUS=200",
		"SCRIPT_FILENAME=" + scriptPath,
	}
	return defaultHandler
}
