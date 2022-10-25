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

type Order struct {
	ID                  int64          `json:"id"`
	ThirdPartyID        int64          `json:"thirdPartyId" binding:"required"`
	OutOrderNo          string         `json:"outOrderNo" binding:"required"`
	OrderStatus         int            `gorm:"default:0" json:"orderStatus"`
	TotalAmount         int64          `json:"totalAmount"`
	TotalDiscountAmount int64          `json:"totalDiscountAmount"`
	RealTotalAmount     int64          `json:"realTotalAmount"`
	PostalFee           int64          `json:"postalFee"`
	Currency            string         `json:"currency"`
	Thumbnail           OrderThumbnail `json:"thumbnail"`

	UnixModelTimeWithDel
}

func (Order) TableName() string {
	return "css_order"
}

type OrderThumbnail []string

type OrderFmtOutPut struct {
	Order
	ThirdPartyOrderFlag string             `json:"thirdPartyOrderFlag"`
	OrderStatusText     string             `json:"orderStatusText"`
	AddressInfo         interface{}        `json:"addressInfo"`
	CreatedTime         string             `json:"createdTime"`
	ProductItems        []*OrderProductFmt `json:"productItems"`
	Extra               interface{}        `json:"extra"`
	OrderHistories      []*OrderHistoryFmt `json:"orderHistories"`
}

type OrderStatistic struct {
	WaitPay     int64 `json:"waitPay"`     // 待下单数量
	WaitGot     int64 `json:"waitGot"`     // 待收货数量
	WaitDeliver int64 `json:"waitDeliver"` // 待发货数量
	// Delivered   int64 `json:"delivered"`   // 已经发货订单数
}

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

func (t OrderThumbnail) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *OrderThumbnail) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), t)
}

func (o Order) CreateBaseOrder(tx *gorm.DB, thirdPartyFlag, outOrderNo string, thumbnails OrderThumbnail, postalFee int64, currency string) (*Order, error) {
	thirdPartyId := o.GetThirdPartyOrderIdByFlag(thirdPartyFlag)
	if thirdPartyId == 0 {
		return nil, fmt.Errorf("invalid third party flag")
	}

	if o.OutOrderNoExistsByThirdPartyId(outOrderNo, thirdPartyId) {
		return nil, fmt.Errorf("out of order number is exists")
	}

	o.ThirdPartyID = thirdPartyId
	o.OutOrderNo = outOrderNo
	if thumbnails != nil {
		o.Thumbnail = thumbnails
	}
	o.PostalFee = postalFee
	o.Currency = currency

	if err := tx.Create(&o).Error; err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create order outOrderNo: %s, error: %v", outOrderNo, err))
	}
	return &o, nil
}

func (o *Order) CalculateAmount(orderProducts []OrderProduct, postalFee int64) error {
	for _, item := range orderProducts {
		fmt.Println(item.TotalAmount)
		o.TotalAmount += item.TotalAmount
		o.TotalDiscountAmount += item.TotalDiscountAmount
		o.RealTotalAmount += item.RealTotalAmount
	}
	o.TotalAmount += postalFee
	o.RealTotalAmount += postalFee
	return nil
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
	case OrderStatusOfCancel:
		return HistoryTypeOfOrderCancel
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
		return append(rt, OrderStatusOfPay, OrderStatusOfCancel, OrderStatusOfGot, OrderStatusOfGotPart, OrderStatusOfDelivery), nil
	case OrderStatusOfPay:
		return append(rt, OrderStatusOfGotPart, OrderStatusOfGot, OrderStatusOfCancel, OrderStatusOfDelivery), nil
	case OrderStatusOfGotPart:
		return append(rt, OrderStatusOfGot, OrderStatusOfCancel, OrderStatusOfDelivery), nil
	case OrderStatusOfGot:
		return append(rt, OrderStatusOfDelivery, OrderStatusOfCancel), nil
	case OrderStatusOfCancel:
		return rt, nil
	default:
		return nil, fmt.Errorf("当前订单状态不可更新")
	}
}

func (o Order) MultiFmt() ([]*OrderFmtOutPut, error) {
	// orderIds := []int64{}
	// 批量获取订单商品
	// opList := OrderProduct{}.Search(nil, []*SearchCond{{ColumnName: "order_id", Operator: "IN (?)", Context: orderIds}}, nil)

	// productIds := []int64{}
	// 批量获取product
	// products := Product{}.Search([]*SearchCond{{ColumnName: "id", Operator: "IN (?)", Context: productIds}})

	// skuIds := []int64{}
	// skus := ProductSku{}.Search([]*SearchCond{{ColumnName: "id", Operator: "IN (?)", Context: skuIds}})
	// 批量获取sku
	// 批量获取订单历史
	return nil, nil

}

func (o Order) Fmt() (*OrderFmtOutPut, error) {
	fmtOrder := OrderFmtOutPut{}
	fmtOrder.CreatedTime = time.Unix(int64(o.CreatedAt), 0).Format("2006-01-02 15:04:05")
	fmtOrder.OrderStatusText = o.GetOrderStatusText(o.OrderStatus)
	fmtOrder.Thumbnail = o.Thumbnail
	fmtOrder.ThirdPartyOrderFlag = o.GetThirdPartyFlag()
	fmtOrder.AddressInfo = o.GetOrderFmtAddress()
	fmtOrder.Extra = o.GetOrderFmtExtend()
	fmtOrder.Order = o

	// 订单商品
	orderProducts := OrderProduct{}.GetByOrderId(o.ID)
	for _, productItem := range orderProducts {
		item := productItem.Fmt()
		fmtOrder.ProductItems = append(fmtOrder.ProductItems, item)
	}

	// 订单历史
	orderHistory := OrderHistory{}.GetByOrderIdDesc(o.ID)
	for _, historyItem := range orderHistory {
		item := historyItem.Fmt()
		fmtOrder.OrderHistories = append(fmtOrder.OrderHistories, item)
	}
	return &fmtOrder, nil
}

func (o Order) Statistic() OrderStatistic {

	rt := OrderStatistic{}

	_ls := []Order{}
	if err := drivers.Mysql().Model(&o).Where("order_status", []int{OrderStatusOfInit, OrderStatusOfPay, OrderStatusOfGot, OrderStatusOfGotPart}).Select("order_status", "id").Find(&_ls).Error; err != nil {
		return rt
	}

	for _, item := range _ls {
		switch item.OrderStatus {
		case OrderStatusOfInit:
			rt.WaitPay += 1
		case OrderStatusOfPay:
			rt.WaitGot += 1
		case OrderStatusOfGot, OrderStatusOfGotPart:
			rt.WaitDeliver += 1
			// case OrderStatusOfDelivery:
			// 	rt.Delivered += 1
		}
	}

	return rt
}
