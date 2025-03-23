package websocket

import "time"

/**
@Author: loser
@Description: the server options of websocket server(注意这一种配置方式的技巧)
*/

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	ack               AckType       // ACK 模式
	ackTimeout        time.Duration // ACK 超时时间
	pattern           string
	maxConnectionIdle time.Duration
	concurrency       int
}

// @Description: apply the options
func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication:    new(authentication),
		pattern:           "/ws",
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultackTimeOut,
		concurrency:       defaultConcurrency,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

func WithServerPattern(pattern string) ServerOptions {
	return func(opt *serverOption) {
		opt.pattern = pattern
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}
