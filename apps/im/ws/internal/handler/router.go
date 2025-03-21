package handler

import (
	"im-chat/apps/im/ws/internal/handler/conversation"
	"im-chat/apps/im/ws/internal/handler/push"
	"im-chat/apps/im/ws/internal/handler/user"
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/apps/im/ws/websocket"
)

/**
@Author: loser
@Description: add routes into websocket server
*/

func RegisterHandlers(srv *websocket.Server, svc *svc.ServerContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(svc),
		},
	})

	srv.AddRoutes([]websocket.Route{
		{
			Method:  "conversation.chat",
			Handler: conversation.Chat(svc),
		},
	})

	srv.AddRoutes([]websocket.Route{
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})
}
