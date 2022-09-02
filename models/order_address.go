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

type OrderAddress struct {
	ID        int64  `json:"id"`
	OrderId   int64  `json:"orderId" binding:"required"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	PostCode  string `json:"postCode"`
	Phone     string `json:"phone"`
	Consignee string `json:"consignee"`
	Extra     string `json:"extra"`
	UnixModelTimeWithDel
}

func (OrderAddress) TableName() string {
	return "css_order_address"
}

type AmzAddress struct {
	Address string `json:"address"`
}

func (OrderAddress) GetByOrderId(orderId int64) (*OrderAddress, error) {
	one := &OrderAddress{}
	if err := drivers.Mysql().Model(one).Where("order_id = ?", orderId).First(&one).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return one, nil
}

func (o OrderAddress) RenderAmzOrderAddress() interface{} {
	var f AmzAddress
	json.Unmarshal([]byte(o.Extra), &f)
	return f.Address
}

func (o *OrderAddress) FmtAmzOrderAddress(address string) error {
	amzAddress := AmzAddress{
		Address: address,
	}

	addr, err := json.Marshal(amzAddress)
	if err != nil {
		return err
	}
	o.Extra = string(addr)
	return nil
}
