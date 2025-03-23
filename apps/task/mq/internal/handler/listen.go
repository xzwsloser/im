package handler

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"im-chat/apps/task/mq/internal/handler/msgTransfer"
	"im-chat/apps/task/mq/internal/svc"
)

type Listen struct {
	svc *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{
		svc: svc,
	}
}

// @Description: load multi consumer which implement the consumer handler
func (l *Listen) Service() []service.Service {
	// 加载多个消费者
	return []service.Service{
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer,
			msgTransfer.NewMsgChatTransfer(l.svc)),
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer,
			msgTransfer.NewMsgReadTransfer(l.svc)),
	}
}
