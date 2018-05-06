package gateway

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"time"
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

func GetEtcdInstancer(prefix string, etcdClient etcdv3.Client, logger log.Logger) *etcdv3.Instancer {
	instancer, err := etcdv3.NewInstancer(etcdClient, prefix, logger)
	if err != nil {
		panic(err)
	}
	return instancer
}
