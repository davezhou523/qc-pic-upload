package po

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"qc/order/internal/logic/po"
	"qc/order/internal/svc"
	"qc/order/internal/types"
)

func PrintedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderAddPrintedRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := po.NewPrintedLogic(r.Context(), svcCtx)
		resp, err := l.Printed(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
