package user

import (
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/apps/im/ws/websocket"
)

/**
@Author: loser
@Description: user handler
*/

func Online(svc *svc.ServerContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, message *websocket.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err: ", err)
	}
}
