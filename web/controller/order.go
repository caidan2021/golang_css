/*
 * @Date: 2022-08-22 20:52:32
 */
package controller

import (
	"fmt"
	"gin/drivers"
	"gin/models"
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
