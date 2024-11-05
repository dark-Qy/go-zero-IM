package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

type (
	// GetUserReq 定义请求和返回格式
	GetUserReq struct {
		Id string `json:"id"`
	}
	GetUserResp struct {
		Id    string
		Name  string
		Phone string
	}
)

type UserServer struct{}

// GetUser 定义一个方法去获取用户信息
func (*UserServer) GetUser(req GetUserReq, resp *GetUserResp) error {
	if u, ok := users[req.Id]; ok {
		*resp = GetUserResp{
			Id:    u.Id,
			Name:  u.Name,
			Phone: u.Phone,
		}
		return nil
	}
	return errors.New("没有找到用户")
}

func main() {
	// 创建好服务
	userServer := new(UserServer)

	// 将服务注册至rpc
	err := rpc.Register(userServer)
	if err != nil {
		log.Fatal("注册服务失败", err)
	}

	// 监听服务请求
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("监听失败", err)
	}
	log.Println("服务启动成功")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接收客户端连接失败", err)
			continue
		}
		// 并发处理客户端请求
		go rpc.ServeConn(conn)
	}

}
