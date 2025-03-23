package push

import (
	"github.com/mitchellh/mapstructure"
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/apps/im/ws/websocket"
	"im-chat/apps/im/ws/ws"
	"im-chat/pkg/constants"
)

func Push(svc *svc.ServerContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))
			return
		}

		switch data.ChatType {
		case constants.SingleChatType:
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			group(srv, &data)
		}
	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 获取发送的目标
	rconn := srv.GetConn(recvId)
	// TODO: 离线状态
	if rconn == nil {
		return nil
	}

	srv.Info("push msg %v ", data)
	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		Msg: ws.Msg{
			MType:   data.MType,
			Content: data.Content,
		},
	}), rconn)
}

func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}

	return nil
}
