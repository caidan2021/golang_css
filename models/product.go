/*
 * @Date: 2022-09-02 19:05:03
 */
package models

import (
	"fmt"
	"gin/drivers"
	"gin/util"
)

const (
	ProductStatusOfInit     = 0
	ProductStatusOfOffLine  = 1
	ProductStatusOfActivity = 2
)

type Product struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Status int64  `json:"status"`
	UnixModelTimeWithDel
}

type ProductFmt struct {
	Product
}

func (Product) TableName() string {
	return "css_product"
}

func (p Product) Search(cond []*SearchCond) []*Product {

	query := drivers.Mysql().Model(&p)
	if len(cond) > 0 {
		for _, v := range cond {
			query = query.Where(fmt.Sprintf("%s %s", v.ColumnName, v.Operator), v.Context)
		}
	}
	_ls := []*Product{}
	if err := query.Find(&_ls).Error; err != nil {
		util.Log.Info("Product search got error, cond:%s, err: %s", cond, err.Error())
		return nil
	}

	if len(_ls) == 0 {
		return nil
	}
	return _ls
}

func (p Product) FindById(id int64) *Product {
	sku := p.Search([]*SearchCond{{ColumnName: "id", Operator: "= ?", Context: id}})
	if sku == nil {
		return nil
	}
	return sku[0]
}

func (p Product) Fmt() *ProductFmt {
	return &ProductFmt{
		Product: p,
	}
}

func (p Product) CanShow() bool {
	return p.Status == ProductStatusOfActivity
}
