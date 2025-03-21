package logic

import (
	"context"
	"github.com/pkg/errors"
	"im-chat/pkg/xerr"

	"im-chat/apps/im/rpc/im"
	"im-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话记录
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// 根据 id 查询消息
	if in.MsgId != "" {
		chatLog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(),
				"find chatlog by msgId err %v , req %v ", err, in.MsgId)
		}

		return &im.GetChatLogResp{
			List: []*im.ChatLog{
				{
					Id:             chatLog.ID.Hex(),
					ConversationId: chatLog.ConversationId,
					SendId:         chatLog.SendId,
					RecvId:         chatLog.RecvId,
					MsgType:        int32(chatLog.MsgType),
					MsgContent:     chatLog.MsgContent,
					ChatType:       int32(chatLog.ChatType),
					SendTime:       chatLog.SendTime,
				},
			},
		}, nil
	}

	// 根据查询时间和会话 Id 查询
	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId,
		in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list by send time err %v , req %v ", err, in)
	}

	res := make([]*im.ChatLog, 0, len(data))
	for _, v := range data {
		res = append(res, &im.ChatLog{
			Id:             v.ID.Hex(),
			ConversationId: v.ConversationId,
			SendId:         v.SendId,
			RecvId:         v.RecvId,
			MsgType:        int32(v.MsgType),
			MsgContent:     v.MsgContent,
			ChatType:       int32(v.ChatType),
			SendTime:       v.SendTime,
		})
	}

	return &im.GetChatLogResp{
		List: res,
	}, nil
}
