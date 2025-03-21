package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"im-chat/apps/task/mq/internal/config"
	"im-chat/apps/task/mq/internal/handler"
	"im-chat/apps/task/mq/internal/svc"
)

var configFile = flag.String("f", "etc/test/task.yaml", "the config file")

func main() {
	// 解析命令行参数
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 注意这里的作用是加载配置文件中的配置信息
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Service() {
		serviceGroup.Add(s)
	}
	fmt.Println("starting mq service at ", ctx.Config.ListenOn)
	serviceGroup.Start() // 开启服务
}
