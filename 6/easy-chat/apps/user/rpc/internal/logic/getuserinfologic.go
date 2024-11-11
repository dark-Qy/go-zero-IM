package logic

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"

	"easy-chat-6/apps/user/rpc/internal/svc"
	"easy-chat-6/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

// 定义用户不存在错误
var (
	ErrUserNotFound = errors.New("用户不存在")
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line
	// 根据用户id在数据库中查找用户实体
	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	var resp user.UserEntity
	copier.Copy(&resp, userEntity)

	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
