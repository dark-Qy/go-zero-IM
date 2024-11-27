package handler

import (
	"easy-chat-6/apps/im/ws/internal/handler/user"
	"easy-chat-6/apps/im/ws/internal/svc"
	"easy-chat-6/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
	})
}
