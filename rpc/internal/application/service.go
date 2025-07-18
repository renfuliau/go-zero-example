// internal/application/order_service.go
package application

import (
	"context"
	"fmt"
	"rpc/internal/domain"
	"rpc/internal/domain/entity"

	"github.com/google/uuid"
)

var _ IService = (*orderServiceImpl)(nil)

type orderServiceImpl struct {
	repo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) IService {
	return &orderServiceImpl{repo: repo}
}

func (s *orderServiceImpl) CreateOrder(ctx context.Context, userID int64, items []*entity.OrderItem) (*entity.Order, error) {
	// 1. 準備領域實體 (Entity)
	order := &entity.Order{
		OrderSN: "SN-" + uuid.New().String(),
		UserID:  userID,
		Status:  entity.StatusPendingPayment,
		Items:   items,
	}
	order.CalculateTotalAmount()

	// 2. 直接呼叫倉儲的 CreateOrder 方法
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	
	return order, nil
}

func (s *orderServiceImpl) GetOrderDetail(ctx context.Context, orderID int64) (*entity.Order, error) {
    return s.repo.FindOrderByID(ctx, orderID)
}

func (s *orderServiceImpl) CancelOrder(ctx context.Context, orderID int64) error {
    order, err := s.repo.FindOrderByID(ctx, orderID)
    if err != nil {
        return err
    }

    if order == nil {
        return fmt.Errorf("order not found")
    }

    // 呼叫領域實體的業務規則
    if !order.CanCancel() {
        return fmt.Errorf("order cannot be cancelled, status is %d", order.Status)
    }

    return s.repo.UpdateOrderStatus(ctx, orderID, entity.StatusCancelled)
}