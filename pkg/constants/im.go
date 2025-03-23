package constants

type MType int

// 消息类型
const (
	TextMType MType = iota
)

type ChatType int

// 聊天会话类型
const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

// 内容类型: 0. 消息  1. 已读未读
type ContentType int

const (
	ContentChatMsg ContentType = iota
	ContentMarkRead
)
