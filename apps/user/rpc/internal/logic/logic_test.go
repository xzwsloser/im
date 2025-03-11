package logic

import (
	"github.com/zeromicro/go-zero/core/conf"
	"im-chat/apps/user/rpc/internal/config"
	"im-chat/apps/user/rpc/internal/svc"
	"path/filepath"
)

/**
@Author: loser
@Description: init the logic obj
*/

var svcCtx *svc.ServiceContext

func init() {
	var c config.Config
	conf.MustLoad(filepath.Join("../../etc/test/user.yaml"), &c)
	svcCtx = svc.NewServiceContext(c)
}
