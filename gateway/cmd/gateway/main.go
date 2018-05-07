package main

import (
	"context"
	"ko/gateway"
)

func main() {
	// 1. 配置
	gateway.InitConfig()

	// 2. 日志系统
	logger := gateway.InitLogger()

	// 3. 服务发现
	var ctx = context.Background()
	etcdClient := gateway.InitEtcd(ctx)

	// 1) 用户中心服务
	router := gateway.InitRouter(logger)
	router.Service("/svc/ucenter", etcdClient)
	router.JwtGet("/svc/ucenter/v1/user/{param}", gateway.GetJwtMiddleware())

	// 2) 订单服务...
	router.Service("/svc/order", etcdClient)
	router.Post("/svc/order/v1/order")

	// 3) xx服务...

	// 4. 启动服务器
	gateway.RunServer(logger, (*gateway.GetConfig())["server_port"], router)
}
