package ws

import "im-chat/pkg/constants"

/**
@Author: loser
@Description: define the struct of message in websocket transport
*/

// mapstructure 标签的作用是指定 any -> map[string]interface{} 类型中的 key 的对应关系
type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
	}

	// 用户发送消息内容
	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		Msg                `mapstructure:"msg"`
		SendTime           int64 `mapstructure:"sendTime"`
	}
)
