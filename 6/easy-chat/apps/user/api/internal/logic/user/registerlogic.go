package user

import (
	"context"
	"easy-chat-6/apps/user/rpc/user"
	"easy-chat-6/pkg/xerr"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"easy-chat-6/apps/user/api/internal/svc"
	"easy-chat-6/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrorUserExist = xerr.New(xerr.SERVER_COMMON_ERROR, "用户已存在")
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	registerLogic, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, errors.WithStack(ErrorUserExist)
	}
	var res types.RegisterResp
	copier.Copy(&res, registerLogic)
	return &res, nil
}
