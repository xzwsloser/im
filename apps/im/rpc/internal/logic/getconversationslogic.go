package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"im-chat/apps/im/immodels"
	"im-chat/pkg/xerr"

	"im-chat/apps/im/rpc/im"
	"im-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// 根据用户查询用户会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}

		return nil, errors.Wrapf(xerr.NewDBErr(), "get conversations failed err %v , req %v ", err, in)
	}

	var res im.GetConversationsResp
	copier.Copy(&res, &data)

	// 根据 id 集合列表查询会话
	ids := make([]string, 0, len(data.ConversationList))
	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}

	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.ListByConversationsIds err %v , req %v ", err, ids)
	}

	for _, conversation := range conversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			// 表示已经存在
			continue
		}

		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			// 存在新的消息
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}

	return &res, nil
}
