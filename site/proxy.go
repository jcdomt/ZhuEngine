package site

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

func (s *Site) SendHttp(rw http.ResponseWriter, req *http.Request) *http.Response {
	// 检查是否是 WebSocket 请求
	if isWebSocketRequest(req) {
		return s.handleWebSocket(rw, req)
	}

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

func isWebSocketRequest(req *http.Request) bool {
	connectionHeader := strings.ToLower(req.Header.Get("Connection"))
	upgradeHeader := strings.ToLower(req.Header.Get("Upgrade"))
	return strings.Contains(connectionHeader, "upgrade") && upgradeHeader == "websocket"
}

func (s *Site) handleWebSocket(rw http.ResponseWriter, req *http.Request) *http.Response {
	// 初始化 WebSocket 连接到后端服务器
	dialer := websocket.Dialer{}
	targetURL := "ws://" + s.Config.Server + req.URL.Path + "?" + req.URL.RawQuery

	// 移除不允许重复的头部字段
	reqHeaders := http.Header{}
	for k, v := range req.Header {
		reqHeaders[k] = v
	}
	reqHeaders.Del("Sec-WebSocket-Extensions")
	reqHeaders.Del("Sec-WebSocket-Version")
	reqHeaders.Del("Sec-WebSocket-Key")
	reqHeaders.Del("upgrade")
	reqHeaders.Del("Connection")
	reqHeaders.Del("Host")
	reqHeaders.Del("Origin")
	//reqHeaders.Add("cookie", req.Header.Get("cookie"))

	backendConn, resp, err := dialer.Dial(targetURL, reqHeaders)
	if err != nil {
		http.Error(rw, "WebSocket proxy error: "+err.Error(), http.StatusBadGateway)
		return nil
	}
	defer backendConn.Close()

	// 初始化与客户端的 WebSocket 连接
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	clientConn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Error(err)
		http.Error(rw, "WebSocket upgrade error: "+err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer clientConn.Close()

	// 启动数据传输协程
	errc := make(chan error, 2)
	go proxyWebSocketConn(clientConn, backendConn, errc)
	go proxyWebSocketConn(backendConn, clientConn, errc)

	// 等待任何一方出错
	// if err := <-errc; err != nil {
	// 	log.Println("WebSocket proxy error:", err)
	// }
	<-errc
	return resp
}

func proxyWebSocketConn(dst, src *websocket.Conn, errc chan error) {
	for {
		messageType, message, err := src.ReadMessage()
		if err != nil {
			errc <- err
			return
		}
		if err := dst.WriteMessage(messageType, message); err != nil {
			errc <- err
			return
		}
	}
}
