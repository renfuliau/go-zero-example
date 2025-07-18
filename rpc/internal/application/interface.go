// IOrderService 應用層服務介面
package application

import (
	"context"
	"rpc/internal/domain/entity"
)

type IService interface {
	CreateOrder(ctx context.Context, userID int64, items []*entity.OrderItem) (*entity.Order, error)
	GetOrderDetail(ctx context.Context, orderID int64) (*entity.Order, error)
	CancelOrder(ctx context.Context, orderID int64) error
}