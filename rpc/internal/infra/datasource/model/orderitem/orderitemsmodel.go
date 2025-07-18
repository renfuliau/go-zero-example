package orderitem

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemsModel = (*customOrderItemsModel)(nil)

type (
	// OrderItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemsModel.
	OrderItemsModel interface {
		orderItemsModel
		FindAllByOrderID(ctx context.Context, orderId int64) ([]*OrderItems, error)
	}

	customOrderItemsModel struct {
		*defaultOrderItemsModel
	}
)

// NewOrderItemsModel returns a model for the database table.
func NewOrderItemsModel(conn sqlx.SqlConn, migratePath string, c cache.CacheConf, opts ...cache.Option) OrderItemsModel {
	return &customOrderItemsModel{
		defaultOrderItemsModel: newOrderItemsModel(conn, migratePath, c, opts...),
	}
}

// FindAllByOrderID
func (m *customOrderItemsModel) FindAllByOrderID(ctx context.Context, orderId int64) ([]*OrderItems, error) {
	var resp []*OrderItems

	query := fmt.Sprintf("select %s from %s where `order_id` = ?", orderItemsRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, orderId)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return []*OrderItems{}, nil
	default:
		return nil, err
	}
}
