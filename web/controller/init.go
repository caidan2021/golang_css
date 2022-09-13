/*
 * @Date: 2022-09-09 23:02:37
 */
package controller

import (
	"gin/models"
	"gin/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(ctx *gin.Context) {
	rt := make(map[string]interface{})
	rt["env"] = util.Env()
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"items": rt}))
}

func AdminMenu(ctx *gin.Context) {
	rt := models.AdminMenu{}.GetMenuList()
	ctx.JSON(http.StatusOK, util.SuccessRespPackage(&gin.H{"items": rt}))
}
