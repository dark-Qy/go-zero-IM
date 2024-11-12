package user

import (
	"context"
	"easy-chat-6/apps/user/api/internal/svc"
	"easy-chat-6/apps/user/api/internal/types"
	"easy-chat-6/apps/user/rpc/user"
	"easy-chat-6/pkg/xerr"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// ErrorUserCreate 定义创建用户失败
	ErrorUserCreate = xerr.New(xerr.SERVER_COMMON_ERROR, "用户登录失败")
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic 用户登入
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, errors.WithStack(ErrorUserCreate)
	}
	var res types.LoginResp

	copier.Copy(&res, loginResp)

	return &res, nil
}
