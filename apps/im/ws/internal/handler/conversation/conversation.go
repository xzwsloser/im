package conversation

import (
	"github.com/mitchellh/mapstructure"
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/apps/im/ws/websocket"
	"im-chat/apps/im/ws/ws"
	"im-chat/apps/task/mq/mq"
	"im-chat/pkg/constants"
	"im-chat/pkg/wuid"
	"time"
)

func Chat(svc *svc.ServerContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 私聊实现
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		if data.ConversationId == "" {
			if data.ChatType == constants.SingleChatType {
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			} else {
				data.ConversationId = data.RecvId
			}
		}

		err := svc.MsgChatTransferChatClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixNano(),
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
		})

		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}

func MarkRead(svc *svc.ServerContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 私聊实现
		var data ws.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ChatType:       data.ChatType,
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})

		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
