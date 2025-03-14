// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"go-zero2/sub3/user_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/userList",
				Handler: UserHandler(serverCtx),
			},
		},
	)
}
