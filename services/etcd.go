package services

import (
	"context"
	"github.com/go-kit/kit/sd/etcdv3"
	"time"
	"github.com/go-kit/kit/log"
)

func InitEtcd(ctx context.Context) etcdv3.Client {
	// 1. 配置
	var etcdServer = "localhost:2379" // in the change from v2 to v3, the schema is no longer necessary if connecting directly to an etcd v3 instance

	options := etcdv3.ClientOptions{
		// Path to trusted ca file
		CACert: "",

		// Path to certificate
		Cert: "",

		// Path to private key
		Key: "",

		// Username if required
		Username: "",

		// Password if required
		Password: "",

		// If DialTimeout is 0, it defaults to 3s
		DialTimeout: time.Second * 3,

		// If DialKeepAlive is 0, it defaults to 3s
		DialKeepAlive: time.Second * 3,
	}

	etcdClient, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}
	return etcdClient
}

func RegisterSvc(etcdClient etcdv3.Client, httpLogger log.Logger) *etcdv3.Registrar{
	prefix     := "/svc/ucenter/"  // known at compile time
	instance   := "localhost:4002"       // taken from runtime or platform, somehow
	key        := prefix + instance    // should be globally unique
	value      := "http://" + instance // based on our transport

	registrar := etcdv3.NewRegistrar(etcdClient, etcdv3.Service{
		Key:   key,
		Value: value,
	}, httpLogger)

	// Register our instance.
	registrar.Register()

	return registrar
}