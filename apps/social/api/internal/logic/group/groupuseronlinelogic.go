package group

import (
	"context"
	"im-chat/apps/social/rpc/socialclient"
	"im-chat/pkg/constants"

	"im-chat/apps/social/api/internal/svc"
	"im-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群在线用户
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})

	if err != nil || len(groupUsers.List) == 0 {
		return &types.GroupUserOnlineResp{}, err
	}

	uids := make([]string, 0, len(groupUsers.List))
	for _, group := range groupUsers.List {
		uids = append(uids, group.GroupId)
	}

	onlines, err := l.svcCtx.Hgetall(constants.REDIS_ONLINE_USER)
	if err != nil {
		return nil, nil
	}

	resOnlineList := make(map[string]bool, len(uids))

	for _, s := range uids {
		if _, ok := onlines[s]; ok {
			resOnlineList[s] = true
		} else {
			resOnlineList[s] = false
		}
	}

	return &types.GroupUserOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
