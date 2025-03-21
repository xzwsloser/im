package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/ws/websocket"
	"im-chat/apps/task/mq/internal/svc"
	"im-chat/apps/task/mq/mq"
	"im-chat/pkg/constants"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

// @Description: implement the cosumer handler and listen the queue
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key = ", key, " value = ", value)

	var (
		data mq.MsgChatTransfer
	)

	// 1. 进行数据的转换
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 2. 记录数据
	if err := m.addChatLog(ctx, &data); err != nil {
		return err
	}

	// 3. 推送消息
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
