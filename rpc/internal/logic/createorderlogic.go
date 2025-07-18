package logic

import (
	"context"
	"strconv"

	"rpc/internal/domain/entity"
	"rpc/internal/svc"
	"rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立訂單
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	// 1. 將 pb.OrderItem 轉換為 entity.OrderItem
	var items []*entity.OrderItem
	for _, item := range in.Items {
		items = append(items, &entity.OrderItem{
			ProductID: strconv.FormatInt(item.ProductId, 10),
			Quantity:  int64(item.Quantity),
			Price:     item.Price,
		})
	}

	// 2. 呼叫應用層服務
	entityOrder, err := l.svcCtx.OrderService.CreateOrder(l.ctx, in.UserId, items)
	if err != nil {
		l.Errorf("CreateOrder failed: %v", err)
		return nil, err
	}

	// 3. 回傳 gRPC 回應
	return &order.CreateOrderResp{
		OrderId: entityOrder.OrderSN, // 回傳訂單編號
		Success: true,
	}, nil
}
