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

type OrderProduct struct {
	ID                  int64  `json:"id"`
	OrderId             int64  `json:"orderId" binding:"required"`      // 订单id
	SkuId               int64  `json:"skuId" binding:"required"`        // sku id
	ProductId           int64  `json:"productId" binding:"required"`    // 商品id
	SkuUnitPrice        int64  `json:"skuUnitPrice" binding:"required"` // sku单价
	Count               int64  `json:"count" binding:"required"`        // 购买数量
	Thumbnail           string `json:"thumbnail" binding:"required"`    // 缩略图
	TotalAmount         int64  `json:"totalAmount"`                     // 折扣前金额
	TotalDiscountAmount int64  `json:"totalDiscountAmount"`             // 总折扣金额
	RealTotalAmount     int64  `json:"realTotalAmount"`                 // 实际金额
	UnixModelTimeWithDel
}

type OrderProductFmt struct {
	OrderProduct
	Product            *ProductFmt            `json:"product"`
	ProductSku         *ProductSkuFmt         `json:"sku"`
	OrderProductExtend *OrderProductExtendFmt `json:"orderProductExtends"`
}

func (OrderProduct) TableName() string {
	return "css_order_product"
}

func (op OrderProduct) CreateBaseOrderProduct(tx *gorm.DB, orderId, productId, skuId, count int64, thumbnail string) (*OrderProduct, error) {

	sku := ProductSku{}.FindById(skuId)
	if sku == nil {
		return nil, fmt.Errorf("sku id %d not found", skuId)
	}
	product := Product{}.FindById(sku.ProductId)
	if product == nil {
		return nil, fmt.Errorf("product id %d not found", productId)
	}

	op = OrderProduct{
		OrderId:      orderId,
		SkuId:        skuId,
		ProductId:    productId,
		SkuUnitPrice: sku.UnitPrice,
		Count:        count,
		Thumbnail:    thumbnail,
	}
	op.calculate()

	if err := tx.Create(&op).Error; err != nil {
		return nil, fmt.Errorf("createOrderProduct failed: %v", err)
	}
	return &op, nil
}

func (op *OrderProduct) calculate() {
	op.TotalAmount = op.SkuUnitPrice * op.Count
	op.TotalDiscountAmount = 0
	op.RealTotalAmount = op.TotalAmount - op.TotalDiscountAmount
}

func (op OrderProduct) Search(tx *gorm.DB, cond []*SearchCond, sortCond []*SortCond) []*OrderProduct {
	var _DB *gorm.DB
	if tx == nil {
		_DB = drivers.Mysql()
	} else {
		_DB = tx
	}

	query := _DB.Model(&op)
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
	_ls := []*OrderProduct{}
	if err := query.Find(&_ls).Error; err != nil {
		util.Log.Info("OrderProduct search got error, cond:%s, err: %s", cond, err.Error())
		return nil
	}

	if len(_ls) == 0 {
		return nil
	}
	return _ls
}
func (op OrderProduct) GetByOrderId(orderId int64) []*OrderProduct {
	return op.Search(nil, []*SearchCond{
		{ColumnName: "order_id", Operator: "= (?)", Context: orderId},
	}, nil)
}

func (op OrderProduct) Fmt() *OrderProductFmt {

	rt := OrderProductFmt{}

	rt.OrderProduct = op

	ope := OrderProductExtend{}.FindByOpId(op.ID)
	if ope != nil {
		rt.OrderProductExtend = ope.Fmt()
	}

	p := Product{}.FindById(op.ProductId)
	if p != nil {
		rt.Product = p.Fmt()
	}

	sku := ProductSku{}.FindById(op.SkuId)
	if sku != nil {
		rt.ProductSku = sku.Fmt()
	}

	return &rt

}
