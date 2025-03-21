package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im-chat/apps/user/models"
	"im-chat/apps/user/rpc/internal/config"
	"im-chat/pkg/constants"
	"im-chat/pkg/ctxdata"
	"time"
)

type ServiceContext struct {
	Config    config.Config
	UserModel models.UsersModel
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		Redis:     redis.MustNewRedis(c.Redisx),
		UserModel: models.NewUsersModel(sqlConn, c.Cache),
	}
}

func (svc *ServiceContext) SetRootToken() error {
	// 1. 生成 jwt
	systemToken, err := ctxdata.GetJwtToken(svc.Config.Jwt.AccessSecret, time.Now().Unix(),
		999999999, constants.SYSTEM_ROOT_UID)
	if err != nil {
		return err
	}
	// 2. 存储到 Redis 中
	err = svc.Redis.Set(constants.REDIS_SYSTEM_ROOT_TOKEN, systemToken)
	if err != nil {
		return err
	}
	return nil
}
