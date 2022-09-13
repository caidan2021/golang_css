/*
 * @Date: 2022-08-18 11:20:46
 */
package routes

import (
	"gin/util"
	"gin/web/controller"
	"gin/web/middlewares"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRegisterRoute(c *gin.Engine) {

	// 加载资源
	loadResource()

	c.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, util.SuccessRespPackage("pong"))
		return
	})

	// 用户登陆
	c.POST("/admin/user/login", controller.LoginAdmin)

	admin := c.Group("/admin").Use(middlewares.AdminAuth())
	{
		// 创建用户
		admin.POST("/user/add", controller.AddAdminUser)
		admin.GET("/user/current", controller.Current)

		admin.GET("/common/mate/index", controller.CommonMate)
		admin.GET("/init", controller.Init)
		admin.GET("/menu", controller.AdminMenu)

		// 订单
		admin.GET("/order/list", controller.OrderList)
		admin.POST("/order/create", controller.AdminCreateOrder)
		admin.POST("/order/change/status", controller.ChangeOrderStatus)
		admin.POST("/order/edit/extra", controller.EditOrderExtra)
	}

}

func loadResource() {
	// 加载静态资源
	loadStatics()
	// 加载html模板
	loadView("resource/view/admin")
}

func loadView(pwd string) {

	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		return
	}

	loaded := false
	for _, fileInfo := range fileInfoList {

		if fileInfo.IsDir() {
			loadView(pwd + "/" + fileInfo.Name())
		} else {
			if !loaded {
				c.LoadHTMLGlob(pwd + "/*.html")
				loaded = true
			}
		}
	}
}

func loadStatics() {
	c.Static("statics", "./resource/statics")
}
