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

type OrderProductExtend struct {
	ID             int64 `json:"id"`
	OrderProductId int64 `json:"orderProductId" binding:"required"` //
	OrderId        int64 `json:"orderId" binding:"required"`        // 订单id
	SkuId          int64 `json:"skuId" binding:"required"`          // sku id
	UnixModelTimeWithDel
}
type OrderProductExtendFmt struct {
	OrderProductExtend
}

func (OrderProductExtend) TableName() string {
	return "css_order_product_extend"
}

func (ope OrderProductExtend) Search(tx *gorm.DB, cond []*SearchCond, sortCond []*SortCond) []*OrderProductExtend {
	var _DB *gorm.DB
	if tx == nil {
		_DB = drivers.Mysql()
	} else {
		_DB = tx
	}

	query := _DB.Model(&ope)
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
	_ls := []*OrderProductExtend{}
	if err := query.Find(&_ls).Error; err != nil {
		util.Log.Info("OrderProductExtend search got error, cond:%s, err: %s", cond, err.Error())
		return nil
	}

	if len(_ls) == 0 {
		return nil
	}
	return _ls
}
func (ope OrderProductExtend) FindByOpId(opId int64) *OrderProductExtend {
	cond := []*SearchCond{
		{ColumnName: "order_product_id", Operator: "= (?)", Context: opId},
	}
	ls := OrderProductExtend{}.Search(nil, cond, nil)
	if ls != nil {
		return ls[0]
	}
	return nil
}

func (ope OrderProductExtend) Fmt() *OrderProductExtendFmt {
	return &OrderProductExtendFmt{
		OrderProductExtend: ope,
	}

}
