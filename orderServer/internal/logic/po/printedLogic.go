package po

import (
	"context"
	"github.com/redis/go-redis/v9"
	"qc/common/globalkey"
	"qc/common/xerr"
	"time"

	"qc/orderServer/internal/svc"
	"qc/orderServer/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrintedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrintedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrintedLogic {
	return &PrintedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PrintedLogic) Printed(req *types.OrderAddPrintedRequest) (resp *types.AddResponse, err error) {
	redisClient := l.svcCtx.RedisClient
	key := globalkey.CacheQcPrinted + "." + time.Now().Format("2006-01-02")
	var redisValue redis.Z
	redisValue.Member = req.Barcode
	redisValue.Score = float64(time.Now().UnixNano())
	_, err = redisClient.ZAdd(l.ctx, key, redisValue).Result()
	if err != nil {
		l.Logger.Errorf("添加已打印条码:%v", err)
		return nil, err
	}
	l.Logger.Info("添加已打印条码:" + req.Barcode)
	orderLogic := NewOrderAddLogic(l.ctx, l.svcCtx)
	_, err = orderLogic.OrderDel(req.Barcode)
	if err != nil {
		return nil, err
	}
	return &types.AddResponse{
		Code: xerr.OK,
		Msg:  xerr.MapErrMsg(xerr.OK),
	}, nil
}
