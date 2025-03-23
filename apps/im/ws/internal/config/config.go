package config

import (
	"github.com/zeromicro/go-zero/core/service"
)

/**
@Author: loser
@Description: the websocket configuration
*/

type Config struct {
	service.ServiceConf
	ListenOn string

	JwtAuth struct {
		AccessSecret string
	}

	Mongo struct {
		Url string
		Db  string
	}

	MsgChatTransfer struct {
		Topic string
		Addrs []string
	}

	MsgReadTransfer struct {
		Topic string
		Addrs []string
	}
}
