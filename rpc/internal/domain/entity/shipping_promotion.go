package entity

// ShippingPromotion 物流活動
type ShippingPromotion struct {
	Name     string
	Discount float64 // 例如 0.1 代表 10% 折扣
	Enabled  bool
}
