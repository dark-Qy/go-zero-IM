package main

import (
	"easy-chat-6/apps/im/ws/internal/config"
	"easy-chat-6/apps/im/ws/internal/handler"
	"easy-chat-6/apps/im/ws/internal/svc"
	"easy-chat-6/apps/im/ws/websocket"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	srv := websocket.NewServer(c.ListenOn)
	defer srv.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("Starting websocket server at %v ...\n", c.ListenOn)
	srv.Start()
}
