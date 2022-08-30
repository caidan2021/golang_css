/*
 * @Date: 2022-08-30 10:43:08
 */
package service

import (
	"encoding/json"
	"fmt"
	"gin/models"

	"gorm.io/gorm"
)

type OrderCreateItem struct {
	ThirdPartyFlag string                 `json:"thirdPartyFlag" bing:"required"`
	OutOrderNo     string                 `json:"outOrderNo" bing:"required"`
	Thumbnails     *models.OrderThumbnail `json:"thumbnails"`
	AddressInfo    string                 `json:"addressInfo"`
	Extra          string                 `json:"extra" bing:"required"`
}

func CreateOrder(tx *gorm.DB, item OrderCreateItem) (*models.Order, error) {
	newOrder := models.Order{}

	thirdPartyId := newOrder.GetThirdPartyOrderIdByFlag(item.ThirdPartyFlag)
	if thirdPartyId == 0 {
		return nil, fmt.Errorf("invalid third party flag")
	}

	if newOrder.OutOrderNoExistsByThirdPartyId(item.OutOrderNo, thirdPartyId) {
		return nil, fmt.Errorf("out of order number is exists")
	}

	addressOK, err := newOrder.AddressCheck(item.AddressInfo)
	if err != nil || !addressOK {
		return nil, fmt.Errorf("address check failed")
	}

	addressJson, err := json.Marshal(item.AddressInfo)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to parse address info, origin data: %s : %v", item.AddressInfo, err))
	}
	extraJson, err := json.Marshal(item.Extra)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to parse extra info, origin data: %s : %v", item.Extra, err))
	}

	newOrder.ThirdPartyID = thirdPartyId
	newOrder.OutOrderNo = item.OutOrderNo
	if item.Thumbnails != nil {
		newOrder.Thumbnail = *item.Thumbnails
	}
	newOrder.AddressInfo = string(addressJson)
	newOrder.Extra = string(extraJson)
	if err := tx.Create(&newOrder).Error; err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create order outOrderNo: %s, error: %v", item.OutOrderNo, err))
	}
	return &newOrder, nil
}
