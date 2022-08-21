/*
 * @Date: 2022-08-19 11:24:36
 */
package middlewares

import (
	"gin/models"
	"gin/util"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if util.IsLocalEnv() {
			ctx.Next()
			return
		}
		authInfo := &util.BaseAuth{
			Name:     ctx.Request.Header.Get("Css-Name"),
			Password: ctx.Request.Header.Get("Css-Password"),
		}
		if authInfo.Name == "" || authInfo.Password == "" {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("登陆失败"))
			return
		}

		admin, _ := models.GetAdminUserByNameAndPw(authInfo.Name, authInfo.Password)
		if admin == nil {
			ctx.AbortWithStatusJSON(403, util.FailedRespPackage("账号密码错误"))
			return
		}
		ctx.Set(util.AdminUserKey, authInfo)
		ctx.Next()
	}
}
