/*
 * @Date: 2022-08-19 11:31:25
 */
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin/drivers"

	"gorm.io/gorm"
)

type OrderExtend struct {
	ID      int64  `json:"id"`
	OrderId int64  `json:"orderId" binding:"required"`
	Extra   string `json:"extra"`
	UnixModelTimeWithDel
}

func (OrderExtend) TableName() string {
	return "css_order_extend"
}

type ExtendFmtItem struct {
	Name string `json:"name" binding:"required"`
	Item string `json:"item" binding:"required"`
}

func (oe OrderExtend) CreateOrderExtend(tx *gorm.DB, order Order, items []ExtendFmtItem) (*OrderExtend, error) {
	oe = OrderExtend{
		OrderId: order.ID,
	}
	if order.IsAmzOrder() {
		if items == nil {
			return nil, nil
		}
		if err := oe.FmtOrderExtend(items); err != nil {
			return nil, fmt.Errorf("createOrderExtend fmt amz order extend failed: %v", err)
		}
		if err := tx.Create(&oe).Error; err != nil {
			return nil, fmt.Errorf("createOrderExtend failed: %v", err)
		}
	}
	return &oe, nil
}

func (OrderExtend) GetByOrderId(orderId int64) (*OrderExtend, error) {
	one := &OrderExtend{}
	if err := drivers.Mysql().Model(one).Where("order_id = ?", orderId).First(&one).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return one, nil
}

func (o *OrderExtend) FmtOrderExtend(extend []ExtendFmtItem) error {
	ext, err := json.Marshal(extend)
	if err != nil {
		return err
	}
	o.Extra = string(ext)
	return nil
}

func (o *OrderExtend) RenderOrderExtend() interface{} {
	var f []ExtendFmtItem
	json.Unmarshal([]byte(o.Extra), &f)
	return f
}
