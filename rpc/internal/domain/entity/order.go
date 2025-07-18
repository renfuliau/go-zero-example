// internal/domain/entity/order.go
package entity

import "time"

// OrderStatus 訂單狀態
type OrderStatus int

const (
	StatusPendingPayment OrderStatus = iota // 0: 待支付
	StatusPaid                              // 1: 已支付
	StatusShipped                           // 2: 已發貨
	StatusCompleted                         // 3: 已完成
	StatusCancelled                         // 4: 已取消
)

// Order 是 "訂單" 聚合的聚合根
type Order struct {
	ID          int64
	OrderSN     string
	UserID      int64
	TotalAmount float64
	Status      OrderStatus
	Items       []*OrderItem // 訂單項目
	CreateTime  time.Time
	UpdateTime  time.Time
}

// OrderItem 是訂單中的一個項目
type OrderItem struct {
	ID         int64
	OrderID    int64
	ProductID  string
	Quantity   int64
	Price      float64
	CreateTime time.Time
	UpdateTime time.Time
}

// CalculateTotalAmount 計算訂單總金額
func (o *Order) CalculateTotalAmount() {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
}

// CanCancel 判斷訂單是否可以被取消
func (o *Order) CanCancel() bool {
	return o.Status == StatusPendingPayment
}
