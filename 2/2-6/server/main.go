package main

import (
	"2-6/proto/user"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
)

/*实现该接口
// 定义服务名
type UserServer interface {
	// 定义rpc方法
	GetUser(context.Context, *GetUserReq) (*GetUserResp, error)
}
*/

type UserServer struct{}

func (u *UserServer) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserResp, error) {
	if u, ok := users[req.Id]; ok {
		return &user.GetUserResp{
			Id:    u.Id,
			Name:  u.Name,
			Phone: u.Phone,
		}, nil
	}
	return nil, errors.New("不存在查询用户")
}

func main() {
	// 创建监听
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("监听失败", err)
	}

	// 创建grpc服务
	s := grpc.NewServer()

	// 注册服务
	user.RegisterUserServer(s, new(UserServer))

	log.Println("服务已启动")

	s.Serve(listen)
}
