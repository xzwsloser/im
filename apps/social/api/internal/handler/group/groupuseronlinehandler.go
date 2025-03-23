package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/apps/social/api/internal/logic/group"
	"im-chat/apps/social/api/internal/svc"
	"im-chat/apps/social/api/internal/types"
)

// 群在线用户
func GroupUserOnlineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupUserOnlineReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupUserOnlineLogic(r.Context(), svcCtx)
		resp, err := l.GroupUserOnline(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
