/*
 * @Date: 2022-08-22 20:52:32
 */
package controller

import (
	"fmt"
	"gin/drivers"
	"gin/models"
	"gin/service"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderList(ctx *gin.Context) {

	ls := []*models.Order{}
	query := drivers.Mysql().Model(&models.Order{})

	page, pageSize := util.Paging(ctx.DefaultQuery("page", "1"), ctx.DefaultQuery("pageSize", "15"))

	orderNo := ctx.Query("orderNo")
	if orderNo != "" {
		query.Where("out_order_no LIKE ?", "%"+orderNo+"%")
	}
	orderStatus := ctx.Query("orderStatus")
	if orderStatus != "" {
		query.Where("order_status = ?", orderStatus)
	}

	var total int64
	query.Count(&total)
	if err := query.Order("id Desc").Offset(page).Limit(pageSize).Find(&ls).Error; err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	res := []*models.OrderFmtOutPut{}
	for _, order := range ls {
		item, err := order.RenderData()
		if err != nil {
			ctx.JSON(http.StatusOK, util.FailedRespPackage(fmt.Sprintf("failed render order :%d, error: %v", order.ID, err)))
			return
		}
		res = append(res, item)
	}

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"list": res, "total": total}))
	return
}

func AdminCreateOrder(ctx *gin.Context) {
	r := service.OrderCreateItem{}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	tx := drivers.Mysql()
	newOrder, err := service.CreateOrder(tx, r)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"orderNo": newOrder.OutOrderNo}))
}

func ChangeOrderStatus(ctx *gin.Context) {
	type changeStatus struct {
		ID          int64 `json:"id" binding:"required"`
		OrderStatus int   `json:"orderStatus" binding:"required"`
	}

	r := changeStatus{}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	orderModel := new(models.Order)
	order, err := orderModel.GetByOrderId(r.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	if order == nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("订单不存在"))
		return
	}

	if ok, err := order.ChangeStatusCheck(r.OrderStatus); err != nil || !ok {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	order.ChangeOrderStatus(r.OrderStatus)
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"id": order.ID, "orderStatus": order.OrderStatus}))
	return
}
