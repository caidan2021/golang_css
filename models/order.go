/*
 * @Date: 2022-08-19 11:31:25
 */
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gin/drivers"
	"gin/util"
	"time"

	"gorm.io/gorm"
)

const (
	OrderStatusOfInit     = 0
	OrderStatusOfPay      = 10
	OrderStatusOfGotPart  = 15
	OrderStatusOfGot      = 20
	OrderStatusOfDelivery = 30

	// text
	OrderStatusOfInitText     = "创建"
	OrderStatusOfPayText      = "已下单"
	OrderStatusOfGotPartText  = "部分到货"
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
	OrderStatusOfGotPart:  OrderStatusOfGotPartText,
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

func (Order) GetOrderStatusText(orderStatus int) string {
	if val, ok := statusToTextMap[int64(orderStatus)]; ok {
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

func (o Order) ChangeStatusCheck(orderStatus int) (bool, error) {
	if _, ok := statusToTextMap[int64(orderStatus)]; !ok {
		return false, fmt.Errorf("状态值不存在")
	}
	changeStatusText := o.GetOrderStatusText(orderStatus)
	if o.OrderStatus == orderStatus {
		return false, fmt.Errorf(fmt.Sprintf("订单已经是【%s】状态", changeStatusText))
	}

	nextStatus, err := o.GetNextStatus(o.OrderStatus)
	if err != nil {
		return false, err
	}
	if ok := util.IsContainInIntArr(nextStatus, orderStatus); !ok {
		return false, fmt.Errorf(fmt.Sprintf("不可更改为【%s】状态", changeStatusText))
	}

	return true, nil
}

func (Order) GetNextStatus(currentOrderStatus int) ([]int, error) {
	var rt []int
	switch currentOrderStatus {
	case OrderStatusOfInit:
		return append(rt, OrderStatusOfPay), nil
	case OrderStatusOfPay:
		return append(rt, OrderStatusOfGotPart, OrderStatusOfGot), nil
	case OrderStatusOfGotPart:
		return append(rt, OrderStatusOfGot), nil
	case OrderStatusOfGot:
		return append(rt, OrderStatusOfDelivery), nil
	default:
		return nil, fmt.Errorf("当前订单状态不可更新")
	}
}

func (o Order) AddressFmt() OrderAddress {
	return OrderAddress{}
}

func (o Order) RenderData() (*OrderFmtOutPut, error) {
	fmtOrder := OrderFmtOutPut{}
	fmtOrder.CreatedTime = time.Unix(int64(o.CreatedAt), 0).Format("2006-01-02 15:00:00")
	fmtOrder.OrderStatusText = o.GetOrderStatusText(o.OrderStatus)
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
