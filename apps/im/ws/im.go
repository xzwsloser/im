package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"im-chat/apps/im/ws/internal/config"
	"im-chat/apps/im/ws/internal/handler"
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/apps/im/ws/websocket"
)

/**
@Author: loser
@Description: the enter point of im service
*/

var configFile = flag.String("f", "etc/test/im.yaml", "the config file")

func main() {
	// 解析命令行参数
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 注意这里的作用是加载配置文件中的配置信息
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServerContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)))
	// websocket.WithMaxConnectionIdle(time.Second*10))

	defer srv.Stop()
	handler.RegisterHandlers(srv, ctx)
	srv.Start()
}
