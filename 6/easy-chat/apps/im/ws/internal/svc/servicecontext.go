package svc

import (
	"easy-chat-6/apps/im/immodels"
	"easy-chat-6/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config
	immodels.ChatLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ChatLogModel: immodels.NewChatLogModel(c.Mongo.Url, c.Mongo.Db, "chat_log"),
	}
}
