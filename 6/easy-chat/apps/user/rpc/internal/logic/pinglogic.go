package logic

import (
	"context"

	"easy-chat-6/apps/user/rpc/internal/svc"
	"easy-chat-6/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 定义用户服务
func (l *PingLogic) Ping(in *user.Request) (*user.Response, error) {
	// todo: add your logic here and delete this line

	return &user.Response{}, nil
}
