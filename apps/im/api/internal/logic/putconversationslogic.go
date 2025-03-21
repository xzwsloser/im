package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"im-chat/apps/im/rpc/imclient"
	"im-chat/pkg/ctxdata"

	"im-chat/apps/im/api/internal/svc"
	"im-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新会话
func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutConversationsLogic) PutConversations(req *types.PutConversationsReq) (*types.PutConversationsResp, error) {
	uid := ctxdata.GetUid(l.ctx)

	var conversationList map[string]*imclient.Conversation

	copier.Copy(&conversationList, req.ConversationList)
	_, err := l.svcCtx.PutConversations(l.ctx, &imclient.PutConversationsReq{
		UserId:           uid,
		ConversationList: conversationList,
	})

	return &types.PutConversationsResp{}, err
}
