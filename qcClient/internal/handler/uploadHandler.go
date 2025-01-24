package handler

import (
	"context"
	"qc/qcClient/internal/logic"
	"qc/qcClient/internal/svc"
)

func UploadHandler(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	l := logic.NewUploadLogic(ctx, svcCtx)
	l.DealUpload()
}
