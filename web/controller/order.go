/*
 * @Date: 2022-08-22 20:52:32
 */
package controller

import (
	"encoding/json"
	"fmt"
	"gin/drivers"
	"gin/models"
	"gin/service"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OrderList(ctx *gin.Context) {

	ls := []*models.Order{}
	query := drivers.Mysql().Model(&models.Order{})

	page, pageSize := util.Paging(ctx.DefaultQuery("page", "1"), ctx.DefaultQuery("pageSize", "10"))

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
		item, err := order.Fmt()
		if err != nil {
			ctx.JSON(http.StatusOK, util.FailedRespPackage(fmt.Sprintf("failed render order :%d, error: %v", order.ID, err)))
			return
		}
		res = append(res, item)
	}

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"list": res, "total": total}))
	return
}

func OrderListStatistics(ctx *gin.Context) {

	result := make(map[string]interface{})

	// 总订单数
	result["totalCount"] = 1000
	// 待下单数量
	result["waitingPlaceCount"] = 100
	// 待收货数量
	result["waitingGotCount"] = 100
	// 已发货数量
	result["deliveredCount"] = 100

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"item": result}))
	return
}

func AdminCreateOrder(ctx *gin.Context) {
	r := service.OrderCreateItem{}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	newOrder, err := service.CreateOrder(r)
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

	order, err := models.Order{}.GetByOrderId(r.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}

	if err := drivers.Mysql().Transaction(func(tx *gorm.DB) error {
		if order == nil {
			return fmt.Errorf("订单不存在")
		}
		if ok, err := order.ChangeStatusCheck(r.OrderStatus); err != nil || !ok {
			return fmt.Errorf(err.Error())
		}
		order.OrderStatus = r.OrderStatus

		if err := tx.Save(&order).Error; err != nil {
			return fmt.Errorf(err.Error())
		}
		// 创建订单历史
		if event := order.GetOrderHistoryEvent(r.OrderStatus); event != "" {
			_, err := models.OrderHistory{}.CreateOrderHistory(tx, order.ID, 1, event, "")
			if err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"id": order.ID, "orderStatus": order.OrderStatus}))
}

func EditOrderExtra(ctx *gin.Context) {
	type orderExtra struct {
		OrderId      int64                   `json:"orderId" binding:"required"`
		AddressExtra string                  `json:"addressExtra"`
		OrderExtra   *[]models.ExtendFmtItem `json:"orderExtra"`
	}
	r := orderExtra{}
	if err := ctx.ShouldBindJSON(&r); err != nil {
		errorMsg := util.ValidatorError(err)
		ctx.JSON(http.StatusOK, util.FailedRespPackage(errorMsg))
		return
	}

	order, err := models.Order{}.GetByOrderId(r.OrderId)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	if err := drivers.Mysql().Transaction(func(tx *gorm.DB) error {
		if r.AddressExtra != "" {
			orderAddress, err := models.OrderAddress{}.GetByOrderId(r.OrderId)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

			if orderAddress == nil {
				newOrderAddress, err := models.OrderAddress{}.CreateOrderAddress(tx, *order, r.AddressExtra)
				if err != nil {
					return err
				}
				orderAddress = newOrderAddress
			}
			if order.IsAmzOrder() {
				if err := orderAddress.FmtAmzOrderAddress(fmt.Sprintf("%v", r.AddressExtra)); err != nil {
					return fmt.Errorf(err.Error())
				}
			}
			if err := tx.Save(&orderAddress).Error; err != nil {
				return fmt.Errorf("update order address failed: %s", err)
			}

		}

		if r.OrderExtra != nil {
			orderExtend, err := models.OrderExtend{}.GetByOrderId(r.OrderId)
			if err != nil {
				return fmt.Errorf(err.Error())
			}

			if orderExtend == nil {
				newOrderExtend, err := models.OrderExtend{}.CreateOrderExtend(tx, *order, *r.OrderExtra)
				if err != nil {
					return err
				}
				orderExtend = newOrderExtend
			}

			ext, err := json.Marshal(&r.OrderExtra)
			if err != nil {
				return err
			}
			orderExtend.Extra = string(ext)
			if err := tx.Save(&orderExtend).Error; err != nil {
				return fmt.Errorf("update order extend failed: %s", err)
			}

		}

		return nil

	}); err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(true))
}
