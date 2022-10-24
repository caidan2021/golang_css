/*
 * @Date: 2022-10-24 14:33:41
 */
package api

import (
	"gin/models"
	"gin/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSkuInfo(ctx *gin.Context) {

	result := make(map[string]interface{})

	skuId := ctx.Query("skuId")
	if skuId == "" {
		ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{}))
		return
	}

	skuIdInt, err := strconv.ParseInt(skuId, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage(err.Error()))
		return
	}
	sku := models.ProductSku{}.FindById(skuIdInt)
	if sku == nil {
		ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{}))
		return
	}
	if !sku.CanShow() {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("sku已经下架"))
		return
	}

	product := models.Product{}.FindById(sku.ProductId)
	if product == nil {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("商品不存在"))
		return
	}
	if !product.CanShow() {
		ctx.JSON(http.StatusOK, util.FailedRespPackage("商品已经下架"))
		return
	}

	result["skuId"] = sku.ID
	result["productId"] = sku.ProductId
	result["productName"] = product.Title
	result["skuName"] = sku.Title
	result["postalFee"] = sku.PostalFee
	result["currency"] = sku.Currency

	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"item": result}))
	return
}
