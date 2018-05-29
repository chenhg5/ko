package main

import (
	"context"
	"net/http"
	"ko/services/ucenter"
	"ko/services/ucenter/middleware"
)

func main() {

	// 1. 配置
	ucenter.InitConfig()

	// 2. 日志系统
	httpLogger := ucenter.InitLogger()

	// 3. 服务发现
	var ctx = context.Background()
	etcdClient := ucenter.InitEtcd(ctx)
	registrar  := ucenter.RegisterSvc(etcdClient, httpLogger)
	// TODO: shutdown空指针报错
	defer registrar.Deregister()

	// 4. 路由服务
	var ucenterSvc ucenter.UcenterServiceInterface
	ucenterSvc = ucenter.UcenterService{}
	ucenterSvc = middleware.InstrumentingMiddleware()(ucenterSvc)

	mux := http.NewServeMux()
	mux.Handle("/svc/ucenter/v1/", ucenter.MakeHandler(ucenterSvc, httpLogger))

	ucenter.RunServer(mux, httpLogger, (*ucenter.GetConfig())["server_port"])
}
