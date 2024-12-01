package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

type Server struct {
	sync.RWMutex // 使用读写锁保证并发安全

	opt            *serverOption
	authentication Authentication // 认证接口

	routes map[string]HandlerFunc
	addr   string
	patten string

	connToUser map[*Conn]string
	userToConn map[string]*Conn

	upgrader websocket.Upgrader // 将http协议升级为websocket协议
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		patten:   opt.patten,
		opt:      &opt,
		upgrader: websocket.Upgrader{},

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		authentication: opt.Authentication,

		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// 如果出现了错误，尝试从错误中恢复过来
	defer func() {
		if err := recover(); err != nil {
			s.Error(err)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}
	// 连接鉴权
	if !s.authentication.Auth(w, r) {
		s.Send(&Message{
			FrameType: FrameData,
			Data:      fmt.Sprintf("不具备访问权限"),
		}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 使用协程运行
	go s.handlerConn(conn)
}

func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 关闭已经登录过的连接
	if c := s.userToConn[uid]; c != nil {
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

// 根据连接对象执行任务处理
func (s *Server) handlerConn(conn *Conn) {
	// 记录连接
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 关闭并删除连接
			s.Close(conn)
			return
		}

		// 请求信息
		var message Message
		json.Unmarshal(msg, &message)

		// 依据请求消息类型分类处理
		switch message.FrameType {
		case FramePing:
			// ping：回复
			s.Send(&Message{FrameType: FramePing}, conn)
		case FrameData:
			// 处理
			if handler, ok := s.routes[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(&Message{
					FrameType: FrameData,
					Data:      fmt.Sprintf("不存在请求方法 %v 请仔细检查", message.Method),
				}, conn)
			}
		}
	}
}

func (s *Server) AddRoutes(routes []Route) {
	for _, route := range routes {
		s.routes[route.Method] = route.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil).Error())
}

func (s *Server) Stop() {
	s.Info("websocket server stop at " + s.addr)
}
