package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"im-chat/apps/user/models"
	"im-chat/apps/user/rpc/internal/svc"
	"im-chat/apps/user/rpc/user"
	"im-chat/pkg/xerr"

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

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		userEntity []*models.Users
		err        error
	)

	if in.Phone != "" {
		users, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "query data from database failed: %v", err)
		}
		userEntity = append(userEntity, users)
	} else if in.Name != "" {
		userEntity, err = l.svcCtx.UserModel.ListByName(l.ctx, in.Name)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "query data form database failed: %v", err)
		}
	} else if len(in.Ids) != 0 {
		userEntity, err = l.svcCtx.UserModel.ListByIds(l.ctx, in.Ids)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "query data from database failed: %v", err)
		}
	}

	var resp []*user.UserEntity
	copier.Copy(&resp, &userEntity)
	return &user.FindUserResp{
		User: resp,
	}, nil
}
