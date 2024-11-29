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

	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn

	upgrader websocket.Upgrader // 将http协议升级为websocket协议
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		patten:   opt.patten,
		upgrader: websocket.Upgrader{},

		connToUser: make(map[*websocket.Conn]string),
		userToConn: make(map[string]*websocket.Conn),

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

	// 将http协议升级为websocket协议
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Error("upgrade:err %v", err)
		return
	}
	// 连接鉴权
	if !s.authentication.Auth(w, r) {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("no access permission")))
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 使用协程运行
	go s.handlerConn(conn)
}

func (s *Server) addConn(conn *websocket.Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*websocket.Conn) []string {

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

func (s *Server) Close(conn *websocket.Conn) {
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

func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
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
func (s *Server) handlerConn(conn *websocket.Conn) {
	for {
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("read message:err %v", err)
			// todo: 关闭连接
			s.Close(conn)
			return
		}

		var message Message
		// 解析消息
		if err := json.Unmarshal(msg, &message); err != nil {
			s.Errorf("unmarshal:err %v", err)
			// todo: 关闭连接
			s.Close(conn)
			return
		}
		// 根据请求方法分发路由
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("method %s not found", message.Method)))
			if err != nil {
				s.Errorf("write message:err %v", err)
				return
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
