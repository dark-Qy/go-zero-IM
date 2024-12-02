package svc

import (
	"easy-chat-6/apps/task/mq/internal/config"
)

type ServiceContext struct {
	config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
