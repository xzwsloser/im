package mq

import "im-chat/pkg/constants"

/**
@Author: loser
@Description: 定义 mq 消息结构
*/

type MsgChatTransfer struct {
	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string `json:"sendId"`
	RecvId             string `json:"recvId"`
	SendTime           int64  `json:"sendTime"`

	constants.MType `json:"mType"`
	Content         string `json:"content"`
}
