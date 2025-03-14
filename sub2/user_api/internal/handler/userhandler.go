package handler

import (
	"go-zero2/common/response"
	"go-zero2/sub2/user_api/internal/logic"
	"go-zero2/sub2/user_api/internal/svc"
	"go-zero2/sub2/user_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUserLogic(r.Context(), svcCtx)
		resp, err := l.User(&req)
		response.Response(r, w, resp, err)
	}
}
