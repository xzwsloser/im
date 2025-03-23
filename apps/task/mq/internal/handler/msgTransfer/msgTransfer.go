package msgTransfer

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"im-chat/apps/im/ws/websocket"
	"im-chat/apps/im/ws/ws"
	"im-chat/apps/social/rpc/socialclient"
	"im-chat/apps/task/mq/internal/svc"
	"im-chat/pkg/constants"
)

type baseMsgTransfer struct {
	svc *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svc:    svc,
		Logger: logx.WithContext(context.Background()),
	}
}

// @Description: 进行消息转发
func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.SingleChatType:
		err = m.single(ctx, data)
	case constants.GroupChatType:
		err = m.group(ctx, data)
	}

	return err

}

func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// RecvId --> GroupId
	// 查找群聊中的用户
	users, err := m.svc.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})

	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))
	for _, member := range users.List {
		if member.UserId == data.SendId {
			continue
		}

		data.RecvIds = append(data.RecvIds, member.UserId)
	}

	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}
