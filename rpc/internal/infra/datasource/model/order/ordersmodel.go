package order

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrdersModel = (*customOrdersModel)(nil)

type (
	// OrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrdersModel.
	OrdersModel interface {
		ordersModel
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

// NewOrdersModel returns a model for the database table.
func NewOrdersModel(conn sqlx.SqlConn, migratePath string, c cache.CacheConf, opts ...cache.Option) OrdersModel {
	return &customOrdersModel{
		defaultOrdersModel: newOrdersModel(conn, migratePath, c, opts...),
	}
}
