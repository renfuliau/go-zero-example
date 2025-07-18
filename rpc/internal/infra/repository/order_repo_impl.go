// internal/infra/repository/order_repo_impl.go
package repository

import (
	"context"
	"fmt"
	"rpc/internal/domain"
	"rpc/internal/domain/entity"
	"rpc/internal/infra/datasource/model/order"
	"rpc/internal/infra/datasource/model/orderitem"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ domain.OrderRepository = (*orderRepositoryImpl)(nil)

type orderRepositoryImpl struct {
	conn            sqlx.SqlConn
	ordersModel     order.OrdersModel
	orderItemsModel orderitem.OrderItemsModel
}

func NewOrderRepository(conn sqlx.SqlConn, ordersModel order.OrdersModel, orderItemsModel orderitem.OrderItemsModel) domain.OrderRepository {
	return &orderRepositoryImpl{
		conn:            conn,
		ordersModel:     ordersModel,
		orderItemsModel: orderItemsModel,
	}
}

// CreateOrder 現在內部處理完整的交易
func (r *orderRepositoryImpl) CreateOrder(ctx context.Context, orderEntity *entity.Order) error {
	// 使用 TransactCtx 開啟交易
	return r.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1. 插入主訂單 (orders)
		// 將 domain entity 轉換為 datasource model
		orderData := &order.Orders{
			OrderSn:     orderEntity.OrderSN,
			UserId:      orderEntity.UserID,
			TotalAmount: orderEntity.TotalAmount,
			Status:      int64(orderEntity.Status),
		}

		// 使用樣板內建的 session 支援
		res, err := r.ordersModel.Insert(ctx, orderData, order.WithSession(session))
		if err != nil {
			return fmt.Errorf("repository: CreateOrder failed during order insert: %w", err)
		}

		newOrderID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("repository: CreateOrder failed during get last insert id: %w", err)
		}

		// 將新產生的 order ID 回填到 entity 中，供後續步驟使用
		orderEntity.ID = newOrderID
		for _, item := range orderEntity.Items {
			item.OrderID = newOrderID
		}

		// 2. 批次插入訂單項目 (order_items)
		// 檢查是否有訂單項目需要插入
		if len(orderEntity.Items) > 0 {
			err = r.createOrderItemsBatch(ctx, session, orderEntity.Items)
			if err != nil {
				return fmt.Errorf("repository: CreateOrder failed during items batch insert: %w", err)
			}
		}

		return nil // 如果一切順利，回傳 nil，交易將會被 COMMIT
	})
}

// createOrderItemsBatch 是一個私有輔助函式，專門用來處理批次插入
func (r *orderRepositoryImpl) createOrderItemsBatch(ctx context.Context, session sqlx.Session, items []*entity.OrderItem) error {
	// --- 批次插入的核心邏輯 ---
	// 1. 準備 SQL 語句
	// 基礎語句是 "INSERT INTO `order_items` (`order_id`, `product_id`, `quantity`, `price`) VALUES "
	query := "INSERT INTO `order_items` (`order_id`, `product_id`, `quantity`, `price`) VALUES "

	// 2. 準備參數
	var valueStrings []string
	var valueArgs []interface{}

	for _, item := range items {
		// 為每一筆資料建立一個 VALUES (?, ?, ?, ?) 子句
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		// 將每一筆資料的欄位值加入到參數列表中
		valueArgs = append(valueArgs, item.OrderID, item.ProductID, item.Quantity, item.Price)
	}

	// 3. 組合完整的 SQL 語句
	// 將所有 VALUES 子句用逗號連接起來，例如：VALUES (?, ?, ?, ?), (?, ?, ?, ?), ...
	stmt := fmt.Sprintf("%s %s", query, strings.Join(valueStrings, ","))

	// 4. 執行批次插入
	// 注意：批次插入無法利用 go-zero 的快取機制，因為快取是針對單筆操作設計的。
	// 但因為這是寫入操作，主要影響的是讀取快取，我們可以在操作後手動清理相關快取。
	// 在這個情境下，因為訂單項目通常是跟著主訂單一起查詢，單獨為它們設快取的效益不大，所以直接寫入即可。
	_, err := session.ExecCtx(ctx, stmt, valueArgs...)
	return err
}

// UpdateOrderStatus 更新訂單狀態
func (r *orderRepositoryImpl) UpdateOrderStatus(ctx context.Context, id int64, status entity.OrderStatus) error {
	orderModel, err := r.ordersModel.FindOne(ctx, id)
	if err != nil {
		return fmt.Errorf("UpdateOrderStatus find order failed: %w", err)
	}
	orderModel.Status = int64(status)
	return r.ordersModel.Update(ctx, orderModel)
}

// FindOrderByID 實作根據 ID 查找訂單
func (r *orderRepositoryImpl) FindOrderByID(ctx context.Context, id int64) (*entity.Order, error) {
	// 讀取操作通常不需要在交易中
	orderModel, err := r.ordersModel.FindOne(ctx, id)
	if err != nil {
		if err == order.ErrNotFound {
			return nil, nil // 找不到時回傳 nil，而不是錯誤
		}
		return nil, fmt.Errorf("repository: FindOrderByID find order failed: %w", err)
	}

	// 查詢關聯的訂單項目
	// 注意：FindAllByOrderID 是一個自訂方法，你需要將它加到 orderitemsmodel.go 中
	itemsModel, err := r.orderItemsModel.FindAllByOrderID(ctx, id)
	if err != nil && err != orderitem.ErrNotFound {
		return nil, fmt.Errorf("repository: FindOrderByID find items failed: %w", err)
	}

	// 將 datasource model 轉換為 domain entity
	return r.toOrderEntity(orderModel, itemsModel), nil
}

// FindOrderByOrderSN 實作根據訂單編號查找訂單
func (r *orderRepositoryImpl) FindOrderByOrderSN(ctx context.Context, orderSN string) (*entity.Order, error) {
	orderModel, err := r.ordersModel.FindOneByOrderSn(ctx, orderSN)
	if err != nil {
		if err == order.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("repository: FindOrderByOrderSN find order failed: %w", err)
	}

	itemsModel, err := r.orderItemsModel.FindAllByOrderID(ctx, orderModel.Id)
	if err != nil && err != orderitem.ErrNotFound {
		return nil, fmt.Errorf("repository: FindOrderByOrderSN find items failed: %w", err)
	}

	return r.toOrderEntity(orderModel, itemsModel), nil
}

// toOrderEntity 是一個輔助函式，用於將資料庫模型轉換為領域實體
func (r *orderRepositoryImpl) toOrderEntity(orderModel *order.Orders, itemsModel []*orderitem.OrderItems) *entity.Order {
	if orderModel == nil {
		return nil
	}

	var itemsEntity []*entity.OrderItem
	for _, item := range itemsModel {
		itemsEntity = append(itemsEntity, &entity.OrderItem{
			ID:         item.Id,
			OrderID:    item.OrderId,
			ProductID:  item.ProductId,
			Quantity:   item.Quantity,
			Price:      item.Price,
			CreateTime: item.CreateTime,
			UpdateTime: item.UpdateTime,
		})
	}

	return &entity.Order{
		ID:          orderModel.Id,
		OrderSN:     orderModel.OrderSn,
		UserID:      orderModel.UserId,
		TotalAmount: orderModel.TotalAmount,
		Status:      entity.OrderStatus(orderModel.Status),
		Items:       itemsEntity,
		CreateTime:  orderModel.CreateTime,
		UpdateTime:  orderModel.UpdateTime,
	}
}
