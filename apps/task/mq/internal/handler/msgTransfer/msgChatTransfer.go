package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im-chat/apps/im/immodels"
	"im-chat/apps/im/ws/ws"
	"im-chat/apps/task/mq/internal/svc"
	"im-chat/apps/task/mq/mq"
	"im-chat/pkg/bitmap"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
	}
}

// @Description: implement the cosumer handler and listen the queue
func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key = ", key, " value = ", value)

	var (
		data  mq.MsgChatTransfer
		msgId = primitive.NewObjectID()
	)

	// 1. 进行数据的转换
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 2. 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		MsgId:          msgId.Hex(),
		SendTime:       data.SendTime,
		MType:          data.MType,
		Content:        data.Content,
	})
}

//func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
//	return m.svc.WsClient.Send(websocket.Message{
//		FrameType: websocket.FrameData,
//		Method:    "push",
//		FromId:    constants.SYSTEM_ROOT_UID,
//		Data:      data,
//	})
//}
//
//func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {
//	// 查询用户,并且推送到消息队列服务
//	users, err := m.svc.GroupUsers(ctx, &socialclient.GroupUsersReq{
//		GroupId: data.RecvId,
//	})
//
//	if err != nil {
//		return err
//	}
//
//	data.RecvIds = make([]string, 0, len(users.List))
//	for _, member := range users.List {
//		if member.UserId == data.SendId {
//			continue
//		}
//
//		data.RecvIds = append(data.RecvIds, member.UserId)
//	}
//
//	return m.svc.WsClient.Send(websocket.Message{
//		FrameType: websocket.FrameData,
//		Method:    "push",
//		FromId:    constants.SYSTEM_ROOT_UID,
//		Data:      data,
//	})
//}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
