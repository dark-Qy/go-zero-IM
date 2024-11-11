package logic

import (
	"context"
	"easy-chat-6/apps/user/models"
	"github.com/jinzhu/copier"

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

// FindUser 查找用户
func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line
	// 可以通过三种方式：电话、名字和ID查找用户
	var (
		userEntitys []*models.Users
		err         error
	)

	// 如果电话不为空，说明通过电话查找
	if in.Phone != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		// 如果名字不为空，说明通过名字查找
		userEntitys, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		// 如果ID不为空，说明通过ID查找
		userEntitys, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, &userEntitys)

	return &user.FindUserResp{
		User: resp,
	}, nil
}
