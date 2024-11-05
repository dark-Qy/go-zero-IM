package main

import (
	"log"
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

func main() {
	// 首先连接服务端
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("建立连接失败", err)
	}
	defer client.Close()

	// 定义请求和返回
	var (
		req  = GetUserReq{Id: "3"}
		resp GetUserResp
	)

	// 调用rpc服务
	err = client.Call("UserServer.GetUser", req, &resp)
	if err != nil {
		log.Println("请求失败", err)
		return
	}
	log.Println(resp)
}
