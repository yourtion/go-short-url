package main

import (
	"github.com/yourtion/go-short-url/internal"
	"github.com/yourtion/go-short-url/internal/services"
)

func main() {

	defer func() {
		// 保证程序退出的时候写入日志
		services.FlushAccessLog()
	}()

	// 启动服务
	internal.StartServer()
}
