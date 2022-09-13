/*
 * @Date: 2022-08-18 11:20:46
 */
package routes

import (
	"gin/util"

	"github.com/gin-gonic/gin"
)

var c *gin.Engine

func HttpRun() {

	c = gin.New()
	// 后台路由注册
	AdminRegisterRoute(c)
	// 视图路由
	ResourceRegisterRoute(c)
	// api路由
	ApiRegisterRoute(c)

	// 启动http服务
	c.Run(util.Config.HttpAddr)
}
