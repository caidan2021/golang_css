/*
 * @Date: 2022-08-19 11:31:25
 */
package models

import (
	"encoding/json"
	"errors"
	"gin/drivers"

	"gorm.io/gorm"
)

type OrderExtend struct {
	ID      int64  `json:"id"`
	OrderId int64  `json:"orderId" binding:"Required"`
	Extra   string `json:"extra"`
	UnixModelTimeWithDel
}

func (OrderExtend) TableName() string {
	return "css_order_extend"
}

type ExtendFmtItem struct {
	Name string `json:"name" binding:"Required"`
	Item string `json:"item" binding:"Required"`
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
