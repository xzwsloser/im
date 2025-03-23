package friend

import (
	"context"
	"im-chat/apps/social/rpc/socialclient"
	"im-chat/pkg/constants"
	"im-chat/pkg/ctxdata"

	"im-chat/apps/social/api/internal/svc"
	"im-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendOnlineResp, err error) {
	uid := ctxdata.GetUid(l.ctx)
	friendList, err := l.svcCtx.FriendList(l.ctx, &socialclient.FriendListReq{
		UserId: uid,
	})

	if err != nil || len(friendList.List) == 0 {
		return &types.FriendOnlineResp{}, nil
	}

	uids := make([]string, 0, len(friendList.List))
	for _, friend := range friendList.List {
		uids = append(uids, friend.FriendUid)
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

	return &types.FriendOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
