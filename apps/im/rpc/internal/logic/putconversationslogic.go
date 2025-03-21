package logic

import (
	"context"
	"github.com/pkg/errors"
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/rpc/im"
	"im-chat/apps/im/rpc/internal/svc"
	"im-chat/pkg/constants"
	"im-chat/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(),
			"ConversationsModels.FindByUserId err %v , req %v ", err, in)
	}

	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	for s, conversation := range in.ConversationList {
		var oldTotal int
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}

		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			Total:          int(conversation.Read) + oldTotal,
			Seq:            conversation.Seq,
		}
	}

	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(),
			"ConversationsModel.Update err %v , req %v ", err, data)
	}

	return &im.PutConversationsResp{}, nil
}
