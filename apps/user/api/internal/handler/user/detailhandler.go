package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"im-chat/apps/user/api/internal/logic/user"
	"im-chat/apps/user/api/internal/svc"
	"im-chat/apps/user/api/internal/types"
)

// 获取用户信息
func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewDetailLogic(r.Context(), svcCtx)
		resp, err := l.Detail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
