/*
 * @Date: 2022-08-19 11:24:36
 */
package middlewares

import (
	"fmt"
	"gin/models"
	"gin/service"
	"gin/util"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 获取token
		token, _ := ctx.Cookie("css-token")
		if token == "" {
			token = ctx.Request.Header.Get("Css-Token")
		}

		if token == "" {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}

		// 解析token
		info, err := service.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}
		if info.ExpiresAt < time.Now().Unix() {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}
		admin, _ := models.GetAdminUserByToken(token)
		if admin != nil {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}
		ctx.Set(util.AdminUserKey, admin)
		fmt.Println(ctx.Get(util.AdminUserKey))
		ctx.Next()
	}
}
