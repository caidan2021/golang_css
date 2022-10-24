/*
 * @Date: 2022-08-30 10:43:08
 */
package service

import (
	"fmt"
	"gin/drivers"
	"gin/models"

	"gorm.io/gorm"
)

type OrderCreateItem struct {
	ThirdPartyFlag string                 `json:"thirdPartyFlag" binding:"required"`
	OutOrderNo     string                 `json:"outOrderNo" binding:"required"`
	Thumbnails     *models.OrderThumbnail `json:"thumbnails"`
	AddressInfo    string                 `json:"addressInfo"`
	PostalFee      int64                  `json:"postalFee" binding:"required"`
	Currency       string                 `json:"currency" binding:"required"`
	ProductItems   []models.OrderProduct  `json:"productItems" binding:"required"`
	Extra          []models.ExtendFmtItem `json:"extra"`
}

func CreateOrder(item OrderCreateItem) (*models.Order, error) {
	newOrder := models.Order{}
	if err := drivers.Mysql().Transaction(func(tx *gorm.DB) error {

		if order, err := newOrder.CreateBaseOrder(tx, item.ThirdPartyFlag, item.OutOrderNo, *item.Thumbnails, item.PostalFee, item.Currency); err != nil {
			return err
		} else {
			newOrder = *order
		}

		// 创建订单商品信息
		if len(item.ProductItems) == 0 {
			return fmt.Errorf("no product items, errorMsg")
		}
		newOps := []models.OrderProduct{}
		for _, productItem := range item.ProductItems {
			op, err := models.OrderProduct{}.CreateBaseOrderProduct(tx, newOrder.ID, productItem.SkuId, productItem.Count, productItem.Thumbnail)
			if err != nil {
				return err
			}
			newOps = append(newOps, *op)
		}

		// 计算订单金额
		if err := newOrder.CalculateAmount(newOps); err != nil {
			return err
		}

		if err := tx.Save(&newOrder).Error; err != nil {
			return fmt.Errorf("createOrderProduct failed: %v", err)
		}
		// 创建订单log
		_, err := models.OrderHistory{}.CreateOrderHistory(tx, newOrder.ID, 1, models.HistoryTypeOfOrderCreate, "")
		if err != nil {
			return err
		}

		// 创建订单扩展
		_, err = models.OrderExtend{}.CreateOrderExtend(tx, newOrder, item.Extra)
		if err != nil {
			return err
		}

		// 创建订单地址信息
		_, err = models.OrderAddress{}.CreateOrderAddress(tx, newOrder, item.AddressInfo)
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &newOrder, nil
}
