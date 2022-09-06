/*
 * @Date: 2022-08-19 11:24:36
 */
package middlewares

import (
	"gin/models"
	"gin/util"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userName := ctx.Request.Header.Get("Css-name")
		password := ctx.Request.Header.Get("Css-password")
		if userName != "" || password != "" {
			if admin, _ := models.GetAdminUserByNameAndPw(userName, password); admin != nil {
				ctx.Set(models.AdminUserKey, admin)
				ctx.Next()
				return
			}
		}

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
		info, err := models.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}
		if info.ExpiresAt < time.Now().Unix() {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}
		admin, _ := models.GetAdminUserByToken(token)

		if admin == nil {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("请登陆"))
			return
		}
		ctx.Set(models.AdminUserKey, admin)
		ctx.Next()
	}
}
