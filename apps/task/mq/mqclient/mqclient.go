package mqclient

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"im-chat/apps/task/mq/mq"
)

/**
@Author: loser
@Description: kafka client
*/

type MsgChatTransferChatClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

type msgChatTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgChatTransferClient(addr []string, topic string,
	opts ...kq.PushOption) MsgChatTransferChatClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.pusher.Push(context.Background(), string(body))
}
