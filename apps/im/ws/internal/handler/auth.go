package handler

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"im-chat/apps/im/ws/internal/svc"
	"im-chat/pkg/ctxdata"
	"net/http"
)

/**
@Author: loser
@Description: websocket 用户鉴权
*/

type JwtAuth struct {
	svc   *svc.ServerContext
	parse *token.TokenParser
	logx.Logger
}

func NewJwtAuth(svc *svc.ServerContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parse:  token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	token, err := j.parse.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		j.Errorf("parse token err %v ", err)
		return false
	}

	if !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.IdentifyKey, claims[ctxdata.IdentifyKey]))
	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	return ctxdata.GetUid(r.Context())
}
