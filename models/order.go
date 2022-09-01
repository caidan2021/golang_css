/*
 * @Date: 2022-08-19 11:31:25
 */
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gin/drivers"
	"time"

	"gorm.io/gorm"
)

const (
	OrderStatusOfInit     = 0
	OrderStatusOfPay      = 10
	OrderStatusOfGot      = 20
	OrderStatusOfDelivery = 30

	// text
	OrderStatusOfInitText     = "创建"
	OrderStatusOfPayText      = "已下单"
	OrderStatusOfGotText      = "到货"
	OrderStatusOfDeliveryText = "已发货"
)

const (
	ThirdPartyOfAmz = 1

	// flag
	ThirdPartyFlagOfAmz = "amz"
)

var orderFlagToTpIdMap map[string]int64 = map[string]int64{
	ThirdPartyFlagOfAmz: ThirdPartyOfAmz,
}
var orderTpIdToFlagMap map[int64]string = map[int64]string{
	ThirdPartyOfAmz: ThirdPartyFlagOfAmz,
}
var statusToTextMap map[int64]string = map[int64]string{
	OrderStatusOfInit:     OrderStatusOfInitText,
	OrderStatusOfPay:      OrderStatusOfPayText,
	OrderStatusOfGot:      OrderStatusOfGotText,
	OrderStatusOfDelivery: OrderStatusOfDeliveryText,
}

type Order struct {
	ID           int64          `json:"id"`
	ThirdPartyID int64          `json:"thirdPartyId" binding:"required"`
	OutOrderNo   string         `json:"outOrderNo" binding:"required"`
	OrderStatus  int            `gorm:"default:0" json:"orderStatus"`
	Thumbnail    OrderThumbnail `json:"thumbnail"`
	AddressInfo  string         `json:"addressInfo"`
	Extra        string         `json:"extra"`
	UnixModelTimeWithDel
}

type OrderThumbnail []string

type OrderFmtOutPut struct {
	Order
	ThirdPartyOrderFlag string       `json:"thirdPartyOrderFlag"`
	OrderStatusText     string       `json:"orderStatusText"`
	AddressInfo         OrderAddress `json:"addressInfo"`
	CreatedTime         string       `json:"createdTime"`
	Extra               string       `json:"extra"`
	OperationBtn        []string     `json:"operationBtn"`
}

type OrderAddress struct {
}

func (Order) TableName() string {
	return "css_order"
}

func (t OrderThumbnail) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *OrderThumbnail) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), t)
}

func (Order) OutOrderNoExistsByThirdPartyId(outOrderNo string, thirdPartyId int64) bool {
	_ls := []*Order{}
	drivers.Mysql().Model(&Order{}).Where("out_order_no = ?", outOrderNo).Where("third_party_id", thirdPartyId).Find(&_ls)
	return len(_ls) > 0
}

func (Order) GetByOutOrderNo(outOrderNo string, thirdPartyId int64) (*Order, error) {
	one := &Order{}
	if err := drivers.Mysql().Model(one).Where("out_order_no = ?", outOrderNo).Where("third_party_id", thirdPartyId).First(&one).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return one, nil
}

func (Order) GetByOrderId(orderId int64) (*Order, error) {
	one := &Order{}
	if err := drivers.Mysql().Model(one).Where("id = ?", orderId).First(&one).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return one, nil
}

func (Order) GetThirdPartyOrderIdByFlag(flag string) int64 {
	if val, ok := orderFlagToTpIdMap[flag]; ok {
		return val
	}
	return 0
}

func (o Order) GetThirdPartyFlag() string {
	if val, ok := orderTpIdToFlagMap[int64(o.ThirdPartyID)]; ok {
		return val
	}
	return "未知"
}

func (o Order) GetOrderStatusText() string {
	if val, ok := statusToTextMap[int64(o.OrderStatus)]; ok {
		return val
	}
	return "未知"
}

func (Order) AddressCheck(address string) (bool, error) {
	return true, nil
	addressTmpl := &OrderAddress{}
	if err := json.Unmarshal([]byte(address), &addressTmpl); err != nil {
		return false, err
	}
	return true, nil
}

func (Order) ChangeStatusCheck(orderStatus int) bool {
	if _, ok := statusToTextMap[int64(orderStatus)]; ok {
		return true
	}
	return false
}

func (o Order) AddressFmt() OrderAddress {
	return OrderAddress{}
}

func (o Order) RenderData() (*OrderFmtOutPut, error) {
	fmtOrder := OrderFmtOutPut{}
	fmtOrder.CreatedTime = time.Unix(int64(o.CreatedAt), 0).Format("2006-01-02 15:00:00")
	fmtOrder.OrderStatusText = o.GetOrderStatusText()
	fmtOrder.Thumbnail = o.Thumbnail
	fmtOrder.ThirdPartyOrderFlag = o.GetThirdPartyFlag()
	fmtOrder.AddressInfo = o.AddressFmt()
	fmtOrder.Order = o
	return &fmtOrder, nil
}

func (o *Order) ChangeOrderStatus(orderStatus int) {
	o.OrderStatus = orderStatus
	drivers.Mysql().Save(&o)
}
