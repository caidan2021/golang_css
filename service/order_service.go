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
	ThirdPartyFlag string                 `json:"thirdPartyFlag" bing:"required"`
	OutOrderNo     string                 `json:"outOrderNo" bing:"required"`
	Thumbnails     *models.OrderThumbnail `json:"thumbnails"`
	AddressInfo    string                 `json:"addressInfo"`
	Extra          []models.ExtendFmtItem `json:"extra"`
}

func CreateOrder(item OrderCreateItem) (*models.Order, error) {
	newOrder := models.Order{}
	if err := drivers.Mysql().Transaction(func(tx *gorm.DB) error {

		thirdPartyId := newOrder.GetThirdPartyOrderIdByFlag(item.ThirdPartyFlag)
		if thirdPartyId == 0 {
			return fmt.Errorf("invalid third party flag")
		}

		if newOrder.OutOrderNoExistsByThirdPartyId(item.OutOrderNo, thirdPartyId) {
			return fmt.Errorf("out of order number is exists")
		}

		newOrder.ThirdPartyID = thirdPartyId
		newOrder.OutOrderNo = item.OutOrderNo
		if item.Thumbnails != nil {
			newOrder.Thumbnail = *item.Thumbnails
		}

		if err := tx.Create(&newOrder).Error; err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to create order outOrderNo: %s, error: %v", item.OutOrderNo, err))
		}

		// 创建订单log

		// 创建订单扩展
		newOrderExtend, err := newOrder.NewOrderExtend(item.Extra)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		if newOrderExtend != nil {
			if err := tx.Create(&newOrderExtend).Error; err != nil {
				return fmt.Errorf("createOrderExtend failed: %v", err)
			}
		}

		// 创建订单地址信息
		newOrderAddress, err := newOrder.NewOrderAddress(item.AddressInfo)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		if newOrderAddress == nil {
			return fmt.Errorf("failed to create order address, errorMsg: %s", item.AddressInfo)
		}
		if err := tx.Create(&newOrderAddress).Error; err != nil {
			return fmt.Errorf("createOrderAddress failed: %v", err)
		}

		// 创建订单商品信息

		return nil

	}); err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return &newOrder, nil
}
