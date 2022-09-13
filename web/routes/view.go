/*
 * @Date: 2022-08-20 17:19:24
 */
package routes

import (
	"gin/web/controller"

	"github.com/gin-gonic/gin"
)

func ResourceRegisterRoute(c *gin.Engine) {
	admin := c.Group("/admin/view/")
	{
		admin.GET("/:pageName", controller.Handler)
	}
}
