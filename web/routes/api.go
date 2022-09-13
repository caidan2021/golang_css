/*
 * @Date: 2022-08-20 17:19:24
 */
package routes

import (
	"gin/web/api"
	"gin/web/middlewares"

	"github.com/gin-gonic/gin"
)

func ApiRegisterRoute(c *gin.Engine) {

	order := c.Group("/order").Use(middlewares.AdminAuth())
	{
		order.POST("/create/batch", api.BatchCreateOrder)
	}

	tools := c.Group("/tools").Use(middlewares.AdminAuth())
	{
		tools.POST("/file/upload", api.Upload)
	}

	img := c.Group("/img")
	{
		img.GET("/:imgName", api.ShowImage)
	}
}
