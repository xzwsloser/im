package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		//AccessExpire int64
	}

	SocialRPC zrpc.RpcClientConf
	UserRPC   zrpc.RpcClientConf
	ImRpc     zrpc.RpcClientConf
	Redisx    redis.RedisConf
}
