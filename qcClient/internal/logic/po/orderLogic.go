package po

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"qc/qcClient/internal/svc"
	"qc/qcClient/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderLogic) Order(req *types.OrderRequest) (resp *types.OrderResponse, err error) {
	// todo: add your logic here and delete this line
	fmt.Println(req.Sno, req.Model, req.Num)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//202302010001CG_TM4C1230H6PMI_leads_
	//202302010001CG_TM4C1230H6PMI_sideView_
	//202302010001CG_TM4C1230H6PMI_topMarking_
	//202302010001CG_TM4C1230H6PMI_undersideMarking_
	//202302010001CG_TM4C1230H6PMI_xRAY_
	barCode := req.Sno + "_" + req.Model + "_leads_"
	fmt.Println(redisClient, barCode)
	//type Z struct {
	//	Score  float64
	//	Member interface{}
	//}
	key := "qcClient.unprint"
	var redisValue redis.Z
	redisValue.Member = barCode
	redisValue.Score = float64(time.Now().UnixMicro())

	//redisValue = new(redis.Z)

	//type Z struct {
	//	Score  float64
	//	Member interface{}
	//}

	res, err := redisClient.ZAdd(l.ctx, key, redisValue).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	//type ZRangeBy struct {
	//	Min, Max      string
	//	Offset, Count int64
	//}
	var zsetRange redis.ZRangeBy
	zsetRange.Min = "-inf"
	zsetRange.Max = "+inf"
	//zsetRange.Offset = 0
	//zsetRange.Count = 10
	data, err := redisClient.ZRangeByScore(l.ctx, key, &zsetRange).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	return &types.OrderResponse{
		Code: 200,
		Msg:  "success",
	}, nil
}
