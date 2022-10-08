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
	OrderStatusOfCancel   = 40

	// text
	OrderStatusOfInitText     = "创建"
	OrderStatusOfPayText      = "已下单"
	OrderStatusOfGotPartText  = "部分到货"
	OrderStatusOfGotText      = "到货"
	OrderStatusOfDeliveryText = "已发货"
	OrderStatusOfCancelText   = "已取消"
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
	OrderStatusOfCancel:   OrderStatusOfCancelText,
}

type Order struct {
	ID           int64          `json:"id"`
	ThirdPartyID int64          `json:"thirdPartyId" binding:"required"`
	OutOrderNo   string         `json:"outOrderNo" binding:"required"`
	OrderStatus  int            `gorm:"default:0" json:"orderStatus"`
	Thumbnail    OrderThumbnail `json:"thumbnail"`

	UnixModelTimeWithDel
}

type OrderThumbnail []string

type OrderFmtOutPut struct {
	Order
	ThirdPartyOrderFlag string          `json:"thirdPartyOrderFlag"`
	OrderStatusText     string          `json:"orderStatusText"`
	AddressInfo         interface{}     `json:"addressInfo"`
	CreatedTime         string          `json:"createdTime"`
	ProductItems        []*OrderProduct `json:"productItems"`
	Extra               interface{}     `json:"extra"`
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

func (Order) GetOrderStatusTextMap() map[int64]string {
	return statusToTextMap
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

func (Order) GetOrderHistoryEvent(orderStatus int) string {
	switch orderStatus {
	case OrderStatusOfPay:
		return HistoryTypeOfOrderDispatch
	case OrderStatusOfGotPart:
		return HistoryTypeOfOrderGotPart
	case OrderStatusOfGot:
		return HistoryTypeOfOrderGot
	case OrderStatusOfDelivery:
		return HistoryTypeOfOrderDelivery
	default:
		return ""
	}

}

func (o Order) GetOrderFmtAddress() interface{} {

	orderAddress, _ := OrderAddress{}.GetByOrderId(o.ID)
	if orderAddress == nil {
		return ""
	}

	// 亚马逊的订单地址
	if o.IsAmzOrder() {
		return orderAddress.RenderAmzOrderAddress()
	}
	return ""
}

func (o Order) GetOrderFmtExtend() interface{} {
	orderExtend, _ := OrderExtend{}.GetByOrderId(o.ID)
	if orderExtend == nil {
		return ""
	}
	return orderExtend.RenderOrderExtend()
}

func (o Order) IsAmzOrder() bool {
	return o.ThirdPartyID == ThirdPartyOfAmz
}

func (o Order) NewOrderAddress(address interface{}) (*OrderAddress, error) {
	newOrderAddress := OrderAddress{
		OrderId: o.ID,
	}
	if o.IsAmzOrder() {
		if err := newOrderAddress.FmtAmzOrderAddress(fmt.Sprintf("%v", address)); err != nil {
			return nil, fmt.Errorf("createOrderAddress fmt amz order address failed: %v", err)
		}
		return &newOrderAddress, nil
	}
	return nil, fmt.Errorf("can't create order address")
}

func (o Order) NewOrderExtend(extend []ExtendFmtItem) (*OrderExtend, error) {
	newOrderExtend := OrderExtend{
		OrderId: o.ID,
	}
	if o.IsAmzOrder() {
		if extend == nil {
			return nil, nil
		}
		if err := newOrderExtend.FmtOrderExtend(extend); err != nil {
			return nil, fmt.Errorf("createOrderExtend fmt amz order extend failed: %v", err)
		}
		return &newOrderExtend, nil
	}
	return nil, nil

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
		return append(rt, OrderStatusOfPay, OrderStatusOfCancel), nil
	case OrderStatusOfPay:
		return append(rt, OrderStatusOfGotPart, OrderStatusOfGot, OrderStatusOfCancel), nil
	case OrderStatusOfGotPart:
		return append(rt, OrderStatusOfGot, OrderStatusOfCancel), nil
	case OrderStatusOfGot:
		return append(rt, OrderStatusOfDelivery, OrderStatusOfCancel), nil
	case OrderStatusOfCancel:
		return rt, nil
	default:
		return nil, fmt.Errorf("当前订单状态不可更新")
	}
}

func (o Order) RenderData() (*OrderFmtOutPut, error) {
	fmtOrder := OrderFmtOutPut{}
	fmtOrder.CreatedTime = time.Unix(int64(o.CreatedAt), 0).Format("2006-01-02 08:09:10")
	fmtOrder.OrderStatusText = o.GetOrderStatusText(o.OrderStatus)
	fmtOrder.Thumbnail = o.Thumbnail
	fmtOrder.ThirdPartyOrderFlag = o.GetThirdPartyFlag()
	fmtOrder.AddressInfo = o.GetOrderFmtAddress()
	fmtOrder.Extra = o.GetOrderFmtExtend()
	fmtOrder.Order = o
	fmtOrder.ProductItems = OrderProduct{}.GetByOrderId(o.ID)
	return &fmtOrder, nil
}
