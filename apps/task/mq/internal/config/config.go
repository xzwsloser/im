package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

/**
@Author: loser
@Description: the configuration of kafka
*/

type Config struct {
	// go-zero 的默认配置,包含日志信息配置等信息
	service.ServiceConf

	// 监听地址
	ListenOn string

	// 转换任务配置信息
	MsgChatTransfer kq.KqConf
	MsgReadTransfer kq.KqConf

	Redisx redis.RedisConf

	Mongo struct {
		Url string
		Db  string
	}

	Ws struct {
		Host string
	}

	MsgReadHandler struct {
		// 是否开启
		GroupMsgReadHandler int
		// 延长时间
		GroupMsgReadRecordDelayTime int
		// 最大缓存条数
		GroupMsgReadRecordDelayCount int
	}

	SocialRpc zrpc.RpcClientConf
}
