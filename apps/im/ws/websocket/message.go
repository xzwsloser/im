package websocket

/**
@Author: loser
@Description: the data type transfer in the websocket connection
*/

type FrameType uint8

const (
	FrameData FrameType = 0x0 // 数据消息
	FramePing FrameType = 0x1 // 心跳消息
	FrameErr  FrameType = 0x3
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string `json:"method"`
	FromId    string `json:"fromId"`
	Data      any    `json:"data"` // interface{} 类型转换为
}

func NewMessage(fromId string, data any) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    fromId,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err,
	}
}
