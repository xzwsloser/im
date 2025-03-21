package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/ws/websocket"
	"im-chat/apps/task/mq/internal/config"
	"im-chat/pkg/constants"
	"net/http"
)

type ServiceContext struct {
	config.Config
	websocket.Client
	*redis.Redis
	immodels.ChatLogModel
	immodels.ConversationModel
	WsClient websocket.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustNewConversationModel(c.Mongo.Url, c.Mongo.Db),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}

	// 注意到需要把 Jwt Token 放入到 http.Header 中
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host,
		websocket.WithClientHeader(header))

	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
