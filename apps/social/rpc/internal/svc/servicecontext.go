package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im-chat/apps/social/rpc/internal/config"
	"im-chat/apps/social/socialmodels"
)

type ServiceContext struct {
	Config             config.Config
	FriendsModel       socialmodels.FriendsModel
	FriendRequestModel socialmodels.FriendRequestsModel
	GroupMembersModel  socialmodels.GroupMembersModel
	GroupRequestsModel socialmodels.GroupRequestsModel
	GroupsModel        socialmodels.GroupsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:             c,
		FriendsModel:       socialmodels.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestModel: socialmodels.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:        socialmodels.NewGroupsModel(sqlConn, c.Cache),
		GroupRequestsModel: socialmodels.NewGroupRequestsModel(sqlConn, c.Cache),
		GroupMembersModel:  socialmodels.NewGroupMembersModel(sqlConn, c.Cache),
	}
}
