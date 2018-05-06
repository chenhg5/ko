package gateway

import (
	jwt "github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

func GetJwtMiddleware() endpoint.Middleware {
	keys := func(token *jwt.Token) (interface{}, error) {
		// jwt密钥
		return []byte("secret"), nil
	}

	// 一个bug跟etcd冲突
	// https://github.com/coreos/etcd/issues/9357
	return kitjwt.NewParser(keys, jwt.SigningMethodHS256, kitjwt.MapClaimsFactory)
}
