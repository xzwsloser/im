package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/rpc/imclient"
	"im-chat/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
	immodels.ConversationModel
	immodels.ConversationsModel
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ChatLogModel:       immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  immodels.MustNewConversationModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: immodels.MustNewConversationsModel(c.Mongo.Url, c.Mongo.Db),
		Im:                 imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
