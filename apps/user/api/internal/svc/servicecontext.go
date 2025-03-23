package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"im-chat/apps/user/api/internal/config"
	"im-chat/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config
	User   userclient.User
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		Redis:  redis.MustNewRedis(c.Redisx),
	}
}
