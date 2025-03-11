package logic

import (
	"context"
	"errors"
	"im-chat/apps/user/rpc/internal/svc"
	"im-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserNotFind = errors.New("没有找到对应用户")
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

// @Description: find the user  with the id
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	userEntity, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if userEntity == nil {
		return nil, ErrUserNotFind
	}

	userRecord := &user.UserEntity{
		Id:       userEntity.Id,
		Avatar:   userEntity.Avatar,
		Nickname: userEntity.Nickname,
		Phone:    userEntity.Phone,
		Sex:      int32(userEntity.Sex.Int64),
	}
	return &user.GetUserInfoResp{
		User: userRecord,
	}, nil
}
