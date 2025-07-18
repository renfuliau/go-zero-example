package logic

import (
	"context"
	"strconv"

	"rpc/internal/svc"
	"rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取得訂單詳情
func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.GetOrderResp, error) {

	orderID, err := strconv.ParseInt(in.OrderId, 10, 64)
	if err != nil {
		l.Errorf("Invalid order id: %s", in.OrderId)
		return nil, err
	}

	// 呼叫應用服務
	entityOrder, err := l.svcCtx.OrderService.GetOrderDetail(l.ctx, orderID)
	if err != nil {
		l.Errorf("GetOrderDetail failed: %v", err)
		return nil, err
	}
	if entityOrder == nil {
		// 處理找不到訂單的情況
		return &order.GetOrderResp{}, nil
	}

	// 將 entity 轉換為 pb
	var pbItems []*order.OrderItem
	for _, item := range entityOrder.Items {
		productID, err := strconv.ParseInt(item.ProductID, 10, 64)
		if err != nil {
			return nil, err
		}
		pbItems = append(pbItems, &order.OrderItem{
			ProductId: productID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}

	return &order.GetOrderResp{
		OrderId:     entityOrder.OrderSN,
		UserId:      entityOrder.UserID,
		Status:      strconv.Itoa(int(entityOrder.Status)), // 將狀態轉為字串
		Items:       pbItems,
		TotalAmount: entityOrder.TotalAmount,
	}, nil
}
