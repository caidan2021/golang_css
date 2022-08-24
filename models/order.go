/*
 * @Date: 2022-08-19 11:31:25
 */
package models

import (
	"database/sql/driver"
	"encoding/json"
	"gin/drivers"
	"time"
)

const (
	OrderStatusOfInit = 0

	// text
	OrderStatusOfInitText = "创建"
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
	OrderStatusOfInit: OrderStatusOfInitText,
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
	OrderStatusText     string       `json:"order"`
	AddressInfo         OrderAddress `json:"addressInfo"`
	CreatedTime         string       `json:"createdTime"`
	Extra               string       `json:"extra"`
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

func (o Order) AddressFmt() OrderAddress {
	return OrderAddress{}
}

func (o Order) RenderData() (*OrderFmtOutPut, error) {
	fmt := OrderFmtOutPut{}
	fmt.Order = o
	fmt.CreatedTime = time.Unix(int64(o.CreatedAt), 0).Format("2006-01-02 15:04:05")
	fmt.OrderStatusText = o.GetOrderStatusText()
	fmt.ThirdPartyOrderFlag = o.GetThirdPartyFlag()
	fmt.AddressInfo = o.AddressFmt()
	return &fmt, nil
}
