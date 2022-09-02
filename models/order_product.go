/*
 * @Date: 2022-08-31 09:08:03
 */
package models

import (
	"gin/drivers"
)

type OrderProduct struct {
	ID                 int64   `json:"id"`
	OrderId            int64   `json:"orderId" binding:"required"`
	ProductId          int64   `json:"productId" binding:"required"`
	SkuId              int64   `json:"skuId" binding:"required"`
	OriginalTotalPrice int64   `json:"originalTotalPrice"`
	RealTotalPrice     int64   `json:"realTotalPrice"`
	TotalDiscount      int64   `json:"totalDiscount"`
	OriginalUnitPrice  int64   `json:"originalUnitPrice"`
	RealUnitPrice      int64   `json:"realUnitPrice"`
	TotalCount         int64   `json:"totalCount" binding:"required"`
	Extra              *string `json:"extra"`
	UnixModelTimeWithDel
}

func (OrderProduct) TableName() string {
	return "css_order_product"
}

func (OrderProduct) NewOrderProduct(orderId, productId, skuId, totalCount int64) *OrderProduct {
	newOrderProduct := OrderProduct{
		OrderId:    orderId,
		ProductId:  productId,
		SkuId:      skuId,
		TotalCount: totalCount,
	}
	return &newOrderProduct
}

func (OrderProduct) GetByOrderId(orderId int64) []*OrderProduct {
	_ls := []*OrderProduct{}
	if err := drivers.Mysql().Model(&OrderProduct{}).Where("order_id = ?", orderId).Find(&_ls).Error; err != nil {
		return nil
	}
	return _ls

}
