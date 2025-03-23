package constants

// 1. 未处理  2. 处理 3. 拒绝
type HandlerResult int

const (
	NoHandlerResult HandlerResult = iota + 1
	PassHandlerResult
	RefuseHandlerResult
	CancelHandlerResult
)

// 群等级: 1. 创建者 2. 管理者 3. 普通
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1
	ManagerGroupRoleLevel
	AtLargeGroupRoleLeve
)

// 进群方式: 1. 邀请 2. 申请
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)
