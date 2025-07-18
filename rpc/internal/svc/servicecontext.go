package svc

import (
	"rpc/internal/application"
	"rpc/internal/config"
	"rpc/internal/infra/datasource/model/order"
	"rpc/internal/infra/datasource/model/orderitem"
	"rpc/internal/infra/repository" // 引入 repository

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	OrderService application.IService // 提供給 logic 層的應用服務
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 建立資料庫連線
	conn := sqlx.NewMysql(c.DataSource.Write)

	// 建立Model
	orderModel := order.NewOrdersModel(conn, c.MigrationPath, c.Cache)
	orderItemModel := orderitem.NewOrderItemsModel(conn, c.MigrationPath, c.Cache)

	// 初始化 Repository
	orderRepo := repository.NewOrderRepository(conn, orderModel, orderItemModel)

	// 初始化 Application Service
	orderService := application.NewOrderService(orderRepo)

	return &ServiceContext{
		Config:       c,
		OrderService: orderService, // 將初始化好的 Service 注入
	}
}
