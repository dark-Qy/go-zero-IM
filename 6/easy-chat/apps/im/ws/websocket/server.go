package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	addr     string
	upgrader websocket.Upgrader
	logx.Logger
}
