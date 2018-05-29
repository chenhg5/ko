package main

import (
	"context"
	"net/http"
	"ko/services"
	"ko/services/middleware"
)

func main() {

	// 1. 配置
	services.InitConfig()

	// 2. 日志系统
	httpLogger := services.InitLogger()

	// 3. 服务发现
	var ctx = context.Background()
	etcdClient := services.InitEtcd(ctx)
	registrar  := services.RegisterSvc(etcdClient, httpLogger)
	// TODO: shutdown空指针报错
	defer registrar.Deregister()

	// 4. 路由服务
	var ucenterSvc services.UcenterServiceInterface
	ucenterSvc = services.UcenterService{}
	ucenterSvc = middleware.InstrumentingMiddleware()(ucenterSvc)

	mux := http.NewServeMux()
	mux.Handle("/svc/ucenter/v1/", services.MakeHandler(ucenterSvc, httpLogger))

	services.RunServer(mux, httpLogger, (*services.GetConfig())["server_port"])
}
