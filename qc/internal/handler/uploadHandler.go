package handler

import (
	"context"
	"qc/qc/internal/logic"
	"qc/qc/internal/svc"
)

func UploadHandler(svcCtx *svc.ServiceContext) {
	ctx := context.Background()
	l := logic.NewUploadLogic(ctx, svcCtx)
	l.DealUpload()
}
