package server

import (
	"github.com/xiaowei520/well_framework/src/server/http"
)

// Serve 启动server的端口监听
func Serve() error {
	errCh := make(chan error)
	//go serveThrift(errCh)
	go serveHTTP(errCh)

	return <-errCh
}

// 启动HTTP服务
func serveHTTP(errCh chan<- error) {
	errCh <- http.Serve()
}

