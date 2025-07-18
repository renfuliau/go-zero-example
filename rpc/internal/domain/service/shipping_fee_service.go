package service

import (
	"rpc/internal/domain/entity"
)

// ShippingFeeService 負責計算訂單運費的領域服務
type ShippingFeeService struct{}

// CalculateShippingFee 計算運費
// 根據訂單總重量、收貨地址地區、會員等級、物流活動等資訊計算運費
func (s *ShippingFeeService) CalculateShippingFee(order *entity.Order, member *entity.Member, address *entity.Address, promo *entity.ShippingPromotion) float64 {
	// 1. 基本運費：每公斤 60 元
	var totalWeight float64
	for _, item := range order.Items {
		totalWeight += float64(item.Quantity) * 0.5 // 假設每件商品 0.5kg
	}
	baseFee := totalWeight * 60
	if baseFee < 80 {
		baseFee = 80 // 最低運費 80 元
	}

	// 2. 地區加價
	var regionFee float64
	switch address.Region {
	case entity.RegionNorth:
		regionFee = 0
	case entity.RegionCentral:
		regionFee = 20
	case entity.RegionSouth:
		regionFee = 30
	case entity.RegionEast:
		regionFee = 40
	case entity.RegionIsland:
		regionFee = 100
	}

	// 3. 會員等級折扣
	var memberDiscount float64
	switch member.Level {
	case entity.LevelNormal:
		memberDiscount = 0
	case entity.LevelSilver:
		memberDiscount = 0.05 // 5% 折扣
	case entity.LevelGold:
		memberDiscount = 0.10 // 10% 折扣
	case entity.LevelPlatinum:
		memberDiscount = 0.15 // 15% 折扣
	}

	// 4. 物流活動折扣
	var promoDiscount float64
	if promo != nil && promo.Enabled {
		promoDiscount = promo.Discount
	}

	// 5. 計算總運費
	fee := baseFee + regionFee
	fee = fee * (1 - memberDiscount)
	fee = fee * (1 - promoDiscount)
	if fee < 0 {
		fee = 0
	}
	return fee
}
