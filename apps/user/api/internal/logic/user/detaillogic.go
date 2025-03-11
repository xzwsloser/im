package user

import (
	"context"
	"github.com/jinzhu/copier"
	"im-chat/apps/user/rpc/user"
	"im-chat/pkg/ctxdata"

	"im-chat/apps/user/api/internal/svc"
	"im-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line
	uid := ctxdata.GetUid(l.ctx) // go-zero save the jwt in the l.context

	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})

	if err != nil {
		return nil, err
	}

	var res types.User
	copier.Copy(&res, userInfoResp.User)
	return &types.UserInfoResp{Info: res}, nil
}
