package main

import (
	"github.com/yourtion/go-short-url/internal"
	"github.com/yourtion/go-short-url/internal/services"
)

func main() {

	defer func() {
		services.FlushAccessLog()
	}()

	// 启动服务
	internal.StartServer()
}
