/*
 * @Date: 2022-09-06 17:21:17
 */
package controller

import (
	"gin/models"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CommonMate(ctx *gin.Context) {
	rt := make(map[string]interface{})
	rt["orderStatusToTextMap"] = util.RenderTwoDimensionForIntStr(models.Order{}.GetOrderStatusTextMap())
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"items": rt}))
}
