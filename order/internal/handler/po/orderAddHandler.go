package po

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"qc/order/internal/logic/po"
	"qc/order/internal/svc"
	"qc/order/internal/types"
)

func OrderAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderAddRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := po.NewOrderAddLogic(r.Context(), svcCtx)
		resp, err := l.OrderAdd(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
