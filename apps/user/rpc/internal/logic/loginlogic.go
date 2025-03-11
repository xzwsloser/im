package logic

import (
	"context"
	"github.com/pkg/errors"
	_ "github.com/pkg/errors"
	"im-chat/apps/user/models"
	"im-chat/apps/user/rpc/internal/svc"
	"im-chat/apps/user/rpc/user"
	"im-chat/pkg/ctxdata"
	"im-chat/pkg/encrypt"
	"im-chat/pkg/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrNotRegister = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号没有注册")
	ErrWrongPass   = xerr.New(xerr.SERVER_COMMON_ERROR, "密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 1. 首先查询用户
	userEntity, err := l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	if err == models.ErrNotFound {
		// only record the stack of error
		return nil, errors.WithStack(ErrPhoneIsRegister)
	}

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v , req %v", err, in.Phone)
	}

	// 2. 验证密码
	if userEntity.Password.String != "" && userEntity.Password.Valid {
		if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
			return nil, errors.WithStack(ErrWrongPass)
		}
	}

	// 3. 生成 token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ctxdata get jwt token err %v", err)
	}

	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
