/*
 * @Date: 2022-09-02 19:05:50
 */
package models

import (
	"fmt"
	"gin/drivers"
	"gin/util"
)

const (
	ProductSkuStatusOfInit     = 0
	ProductSkuStatusOfOffLine  = 1
	ProductSkuStatusOfActivity = 2
)

type ProductSku struct {
	ID        int64  `json:"id"`
	ProductId int64  `json:"productId"`
	Status    int64  `json:"status"`
	Title     string `json:"title"`
	UnitPrice int64  `json:"unitPrice"` // 单位 分
	PostalFee int64  `json:"postalFee"`
	Currency  string `json:"currency"`
	UnixModelTimeWithDel
}
type ProductSkuFmt struct {
	ProductSku
}

func (ProductSku) TableName() string {
	return "css_product_sku"
}

func (ps ProductSku) Search(cond []*SearchCond) []*ProductSku {

	query := drivers.Mysql().Model(&ps)
	if len(cond) > 0 {
		for _, v := range cond {
			query = query.Where(fmt.Sprintf("%s %s", v.ColumnName, v.Operator), v.Context)
		}
	}
	_ls := []*ProductSku{}
	if err := query.Find(&_ls).Error; err != nil {
		util.Log.Info("ProductSku search got error, cond:%s, err: %s", cond, err.Error())
		return nil
	}

	if len(_ls) == 0 {
		return nil
	}
	return _ls
}
func (ps ProductSku) FindById(id int64) *ProductSku {
	sku := ps.Search([]*SearchCond{{ColumnName: "id", Operator: "= ?", Context: id}})
	if sku == nil {
		return nil
	}
	return sku[0]
}

func (ps ProductSku) Fmt() *ProductSkuFmt {
	return &ProductSkuFmt{
		ProductSku: ps,
	}
}

func (p ProductSku) CanShow() bool {
	return p.Status == ProductSkuStatusOfActivity
}
