package po

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"qc/common/globalkey"
	"qc/order/internal/svc"
	"qc/order/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderListLogic) OrderList(req *types.OrderListReq) (resp *types.OrderListResp, err error) {
	fmt.Println("Category:", req.Category)
	fmt.Println("PageSize:", req.PageSize)
	redisClient := l.svcCtx.RedisClient
	////1:未打印流号号，2当日已打印的流水号
	var key string
	var data []string
	if req.Category == 1 {
		key = globalkey.CacheQcUnprint
		var zsetRange redis.ZRangeBy
		zsetRange.Min = "-inf"
		zsetRange.Max = "+inf"
		zsetRange.Offset = 0
		zsetRange.Count = req.PageSize
		data, err = redisClient.ZRangeByScore(l.ctx, key, &zsetRange).Result()
	} else if req.Category == 2 {
		key = globalkey.CacheQcPrinted + "." + time.Now().Format("2006-01-02")
		var zsetRange redis.ZRangeBy
		zsetRange.Min = "-inf"
		zsetRange.Max = "+inf"
		zsetRange.Offset = 0
		zsetRange.Count = req.PageSize
		data, err = redisClient.ZRevRangeByScore(l.ctx, key, &zsetRange).Result()
	}
	var res types.OrderListResp
	if err != nil {
		l.Logger.Error("获取流水号失败:%v", err)
		return nil, err
	}
	res.List = data
	return &res, nil
}
