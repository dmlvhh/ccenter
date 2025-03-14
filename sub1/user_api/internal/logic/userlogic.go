package logic

import (
	"context"
	"go-zero2/sub1/user_models"

	"go-zero2/sub1/user_api/internal/svc"
	"go-zero2/sub1/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.UserReq) (resp *types.UserResp, err error) {
	var user user_models.SysUser
	err = l.svcCtx.DB.Model(&user_models.SysUser{}).First(&user).Error
	return &types.UserResp{
		UUID:     user.UUID.String(),
		Username: user.Username,
		NickName: user.NickName,
	}, nil
}
