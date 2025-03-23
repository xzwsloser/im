package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	// 使用了这一个配置之后居然自动生成了 JwtAuth 的配置
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	UserRPC zrpc.RpcClientConf

	Redisx redis.RedisConf
}
