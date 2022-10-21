/*
 * @Date: 2022-08-31 09:08:03
 */
package models

import (
	"fmt"
	"gin/drivers"
	"gin/util"

	"gorm.io/gorm"
)

const (
	HistoryTypeOfOrderCreate   = "order_create"
	HistoryTypeOfOrderDispatch = "order_dispatch"
	HistoryTypeOfOrderGotPart  = "order_go_part"
	HistoryTypeOfOrderGot      = "order_got"
	HistoryTypeOfOrderDelivery = "order_delivery"
	HistoryTypeOfOrderCancel   = "order_cancel"
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
type OrderHistoryFmt struct {
	OrderHistory
}

func (OrderHistory) TableName() string {
	return "css_order_history"
}

func (oh OrderHistory) CreateOrderHistory(tx *gorm.DB, orderId, actionUserId int64, event, remark string) (*OrderHistory, error) {
	oh = OrderHistory{
		OrderId:    orderId,
		ActionUser: actionUserId,
		Event:      event,
		Remark:     remark,
	}

	if err := tx.Create(&oh).Error; err != nil {
		return nil, fmt.Errorf("createOrderHistory failed: %s", err)
	}
	return &oh, nil

}

func (oh OrderHistory) Search(tx *gorm.DB, cond []*SearchCond, sortCond []*SortCond) []*OrderHistory {
	var _DB *gorm.DB
	if tx == nil {
		_DB = drivers.Mysql()
	} else {
		_DB = tx
	}

	query := _DB.Model(&oh)
	if len(cond) > 0 {
		for _, v := range cond {
			query = query.Where(fmt.Sprintf("%s %s", v.ColumnName, v.Operator), v.Context)
		}
	}
	if len(sortCond) > 0 {
		for _, v := range sortCond {
			query = query.Order(fmt.Sprintf("%s %s", v.ColumnName, v.Sort))
		}
	}

	_ls := []*OrderHistory{}
	if err := query.Find(&_ls).Error; err != nil {
		util.Log.Info("OrderHistory search got error, cond:%s, err: %s", cond, err.Error())
		return nil
	}

	if len(_ls) == 0 {
		return nil
	}
	return _ls
}

func (oh OrderHistory) GetByOrderIdDesc(orderId int64) []*OrderHistory {
	return oh.Search(nil, []*SearchCond{
		{ColumnName: "order_id", Operator: "= (?)", Context: orderId},
	}, []*SortCond{
		{ColumnName: "id", Sort: "desc"},
	})
}

func (oh OrderHistory) Fmt() *OrderHistoryFmt {
	return &OrderHistoryFmt{
		oh,
	}
}
