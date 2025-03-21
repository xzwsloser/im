package websocket

import "time"

const (
	// ACK 应答超时时间
	defaultackTimeOut        = 30 * time.Second
	defaultMaxConnectionIdle = 10 * time.Second
)
