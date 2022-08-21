/*
 * @Date: 2022-08-19 11:28:43
 */
package web

import (
	"encoding/json"
	"fmt"
	"gin/drivers"
	"gin/models"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createItem struct {
	ThirdPartyFlag string                 `json:"thirdPartyFlag" bing:"required"`
	OutOrderNo     string                 `json:"outOrderNo" bing:"required"`
	Thumbnails     *models.OrderThumbnail `json:"thumbnails"`
	AddressInfo    string                 `json:"addressInfo"`
	Extra          string                 `json:"extra" bing:"required"`
}

type createItemReq struct {
	OutOrderNo string `json:"outOrderNo"`
	Rt         bool   `json:"rt"`
	ErrorMsg   string `json:"errorMsg"`
}

func OrderList(ctx *gin.Context) {
	page, pageSize := util.Paging(ctx.DefaultQuery("page", "1"), ctx.DefaultQuery("pageSize", "15"))
	ls := []*models.Order{}
	query := drivers.Mysql().Model(&models.Order{})

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

func BatchCreateOrder(ctx *gin.Context) {
	type req struct {
		Items []createItem `json:"items" form:"items" binding:"required"`
	}

	r := req{}

	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	res := []*createItemReq{}

	err := drivers.Mysql().Transaction(func(tx *gorm.DB) error {

		for _, item := range r.Items {

			itemRt := createItemReq{
				OutOrderNo: item.OutOrderNo,
				Rt:         true,
			}

			newOrder := models.Order{}

			thirdPartyId := newOrder.GetThirdPartyOrderIdByFlag(item.ThirdPartyFlag)
			if thirdPartyId == 0 {
				itemRt.setCreateItemFailed("invalid third party flag")
				res = append(res, &itemRt)
				continue
			}

			if newOrder.OutOrderNoExistsByThirdPartyId(item.OutOrderNo, thirdPartyId) {
				itemRt.setCreateItemFailed("out order is exists")
				res = append(res, &itemRt)
				continue
			}

			addressOK, err := newOrder.AddressCheck(item.AddressInfo)
			if err != nil || !addressOK {
				itemRt.setCreateItemFailed("invalid address")
				res = append(res, &itemRt)
				continue
			}

			addressJson, err := json.Marshal(item.AddressInfo)
			if err != nil {
				itemRt.setCreateItemFailed(fmt.Sprintf("failed to parse address info, origin data: %s : %v", item.AddressInfo, err))
				res = append(res, &itemRt)
				continue
			}
			extraJson, err := json.Marshal(item.Extra)
			if err != nil {
				itemRt.setCreateItemFailed(fmt.Sprintf("failed to parse extra info, origin data: %s : %v", item.Extra, err))
				res = append(res, &itemRt)
				continue
			}

			newOrder.ThirdPartyID = thirdPartyId
			newOrder.OutOrderNo = item.OutOrderNo
			fmt.Println(item.Thumbnails)
			if item.Thumbnails != nil {
				newOrder.Thumbnail = *item.Thumbnails
			}
			newOrder.AddressInfo = string(addressJson)
			newOrder.Extra = string(extraJson)
			if err := tx.Create(&newOrder).Error; err != nil {
				itemRt.setCreateItemFailed(fmt.Sprintf("failed to create order outOrderNo: %s, error: %v", item.OutOrderNo, err))
				res = append(res, &itemRt)
				continue
			}
			res = append(res, &itemRt)
		}
		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"rt": res}))
	return
}

func (c *createItemReq) setCreateItemFailed(errorMsg string) {
	c.Rt = false
	c.ErrorMsg = errorMsg
}
