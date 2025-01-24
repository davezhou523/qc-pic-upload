package po

import (
	"context"
	"github.com/redis/go-redis/v9"
	"qc/common/globalkey"
	"qc/common/xerr"
	"strings"
	"time"

	"qc/orderServer/internal/svc"
	"qc/orderServer/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var picCategory = []string{
	"leads",
	"sideView",
	"topMarking",
	"undersideMarking",
	"xRAY",
}

type OrderAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderAddLogic {
	return &OrderAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderAddLogic) OrderAdd(req *types.OrderAddRequest) (resp *types.AddResponse, err error) {
	redisClient := l.svcCtx.RedisClient
	//202302010001CG_TM4C1230H6PMI_leads_
	//202302010001CG_TM4C1230H6PMI_sideView_
	//202302010001CG_TM4C1230H6PMI_topMarking_
	//202302010001CG_TM4C1230H6PMI_undersideMarking_
	//202302010001CG_TM4C1230H6PMI_xRAY_
	filter := []string{":", "#", "+", "/", "&", ","}
	for _, v := range picCategory {
		for _, v := range filter {
			//过滤特殊字符
			req.Model = strings.ReplaceAll(req.Model, v, "")
		}
		barcode := req.Sno + "_" + req.Model + "_" + v + "_"
		l.Logger.Info("add barcode:" + barcode)
		key := globalkey.CacheQcUnprint
		var redisValue redis.Z
		redisValue.Member = barcode
		redisValue.Score = float64(time.Now().UnixNano())
		_, err := redisClient.ZAdd(l.ctx, key, redisValue).Result()
		if err != nil {
			l.Logger.Errorf("添加条码:%v", err)
			return nil, err
		}

	}

	return &types.AddResponse{
		Code: xerr.OK,
		Msg:  xerr.MapErrMsg(xerr.OK),
	}, nil
}

func (l *OrderAddLogic) OrderDel(barcode string) (resp *types.AddResponse, err error) {
	redisClient := l.svcCtx.RedisClient
	key := globalkey.CacheQcUnprint
	_, err = redisClient.ZRem(l.ctx, key, barcode).Result()
	if err != nil {
		l.Logger.Errorf("删除条码:%v,%v", barcode, err)
		return nil, err
	}

	return &types.AddResponse{
		Code: xerr.OK,
		Msg:  xerr.MapErrMsg(xerr.OK),
	}, nil
}
