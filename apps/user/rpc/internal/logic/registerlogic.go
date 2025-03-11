package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"im-chat/apps/user/models"
	"im-chat/apps/user/rpc/internal/svc"
	"im-chat/apps/user/rpc/user"
	"im-chat/pkg/ctxdata"
	"im-chat/pkg/encrypt"
	"im-chat/pkg/generator"
	"im-chat/pkg/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号已经被注册了")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 1. 首先查询是否注册
	userEntity, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		return nil, errors.WithStack(err)
	}

	if userEntity != nil {
		return nil, errors.Wrapf(ErrPhoneIsRegister, "find user entity in database ")
	}

	uuid := generator.GeneratorUUid()
	userEntity = &models.Users{
		Id:       uuid,
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 2. 加密密码
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	// 3. 插入数据
	_, err = l.svcCtx.UserModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 4. 生成 jwt token 并且返回
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret,
		now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "create jwt failed: %v", err)
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
