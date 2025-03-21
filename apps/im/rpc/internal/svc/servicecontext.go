package svc

import (
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
	immodels.ConversationModel
	immodels.ConversationsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ChatLogModel:       immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  immodels.MustNewConversationModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: immodels.MustNewConversationsModel(c.Mongo.Url, c.Mongo.Db),
	}
}
