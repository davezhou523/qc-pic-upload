// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	po "qc/orderServer/internal/handler/po"
	"qc/orderServer/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/orderServer/add",
				Handler: po.OrderAddHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/orderServer/getList",
				Handler: po.OrderListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/orderServer/addPrinted",
				Handler: po.PrintedHandler(serverCtx),
			},
		},
	)
}
