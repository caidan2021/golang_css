/*
 * @Date: 2022-08-19 11:28:43
 */
package web

import (
	"gin/service"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateItemRes struct {
	OutOrderNo string `json:"outOrderNo"`
	Rt         bool   `json:"rt"`
	ErrorMsg   string `json:"errorMsg"`
}

func BatchCreateOrder(ctx *gin.Context) {
	type req struct {
		Items []service.OrderCreateItem `json:"items" form:"items" binding:"required"`
	}

	r := req{}

	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	res := []*CreateItemRes{}
	for _, item := range r.Items {
		itemRt := CreateItemRes{
			OutOrderNo: item.OutOrderNo,
			Rt:         true,
		}
		if _, err := service.CreateOrder(item); err != nil {
			itemRt.setCreateItemFailed(err.Error())
		}
		res = append(res, &itemRt)
	}

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"rt": res}))
	return
}

func (c *CreateItemRes) setCreateItemFailed(errorMsg string) {
	c.Rt = false
	c.ErrorMsg = errorMsg
}
