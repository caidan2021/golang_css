/*
 * @Date: 2022-08-31 09:08:03
 */
package models

const (
	HistoryTypeOfOrderCreate   = "order_create"
	HistoryTypeOfOrderDispatch = "order_dispatch"
	HistoryTypeOfOrderGotPart  = "order_go_part"
	HistoryTypeOfOrderGot      = "order_got"
	HistoryTypeOfOrderDelivery = "order_delivery"
)

type OrderHistory struct {
	ID         int64   `json:"id"`
	OrderId    int64   `json:"orderId" binding:"required"`
	Event      string  `json:"event" binding:"required"`
	ActionUser int64   `json:"actionUser" binding:"required"`
	Remark     string  `json:"remark"`
	Extra      *string `json:"extra"`

	UnixModelTimeWithDel
}

func (OrderHistory) TableName() string {
	return "css_order_history"
}

func (OrderHistory) NewOrderHistory(orderId, actionUserId int64, event, remark string) *OrderHistory {
	newOrderHistory := OrderHistory{
		OrderId:    orderId,
		ActionUser: actionUserId,
		Event:      event,
		Remark:     remark,
	}
	return &newOrderHistory
}
