package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type Server struct {
	addr     string
	upgrader websocket.Upgrader // 将http协议升级为websocket协议
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		upgrader: websocket.Upgrader{},
		Logger:   logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil).Error())
}

func (s *Server) Stop() {
	s.Info("websocket server stop at " + s.addr)
}
