// internal/domain/repository.go
package domain

import (
	"context"
	"rpc/internal/domain/entity"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OptionFunc 和 Option 保持不變，它們是實現這個模式的關鍵
type OptionFunc func(*Option)

type Option struct {
	Session sqlx.Session
}

func WithSession(session sqlx.Session) OptionFunc {
	return func(o *Option) {
		o.Session = session
	}
}

// OrderRepository 定義了訂單的持久化操作介面
type OrderRepository interface {
	// 負責完整地建立一個訂單聚合 (包含交易和批次插入)
	CreateOrder(ctx context.Context, order *entity.Order) error

	// 更新訂單狀態
	UpdateOrderStatus(ctx context.Context, id int64, status entity.OrderStatus) error

	// 查詢完整的訂單聚合
	FindOrderByID(ctx context.Context, id int64) (*entity.Order, error)
	FindOrderByOrderSN(ctx context.Context, orderSN string) (*entity.Order, error)
}
