package logic

import (
	"context"
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/apps/im/ws/websocket"
	"im-chat/apps/im/ws/ws"
	"im-chat/pkg/constants"
	"im-chat/pkg/wuid"
	"time"
)

type Conversation struct {
	ctx context.Context
	srv *websocket.Server
	svc *svc.ServerContext
}

func NewConversation(ctx context.Context, srv *websocket.Server, svc *svc.ServerContext) *Conversation {
	return &Conversation{
		ctx: ctx,
		srv: srv,
		svc: svc,
	}
}

func (l *Conversation) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       constants.SingleChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       time.Now().UnixNano(),
	}

	err := l.svc.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
