package svc

import (
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/ws/internal/config"
	"im-chat/apps/task/mq/mqclient"
)

/**
@Author: loser
@Description: the service context in im
*/

type ServerContext struct {
	Config config.Config

	immodels.ChatLogModel
	mqclient.MsgChatTransferChatClient
}

func NewServerContext(c config.Config) *ServerContext {
	return &ServerContext{
		Config: c,
		// 注意 MongoDB 中每一个用户操作一个 collection
		ChatLogModel: immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		MsgChatTransferChatClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs,
			c.MsgChatTransfer.Topic),
	}
}
