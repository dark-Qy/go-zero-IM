package logic

import (
	"context"

	"easy-chat-6/apps/user/rpc/internal/svc"
	"easy-chat-6/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查找用户
func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line

	return &user.FindUserResp{}, nil
}
