package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"im-chat/apps/social/rpc/internal/svc"
	"im-chat/apps/social/rpc/social"
	"im-chat/apps/social/socialmodels"
	"im-chat/pkg/constants"
	"im-chat/pkg/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// 1. 判断申请人和目标是否是好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend by uid and fid failed  err: %v req: %v", err, in)
	}

	if friends != nil {
		return &social.FriendPutInResp{}, nil
	}

	// 2. 查看是否申请过
	friendReqs, err := l.svcCtx.FriendRequestModel.FindByReqIdAndUid(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendrequest by rid and uid failed  err: %v req: %v", err, in)
	}

	if friendReqs != nil {
		return &social.FriendPutInResp{}, nil
	}

	// 3. 创建申请记录, 注意需要使用事务
	_, err = l.svcCtx.FriendRequestModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err: %v , request: %v", err, in)
	}

	return &social.FriendPutInResp{}, nil
}
