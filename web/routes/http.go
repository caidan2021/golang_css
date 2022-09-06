/*
 * @Date: 2022-08-18 11:20:46
 */
package routes

import (
	"gin/util"
	"gin/web"
	"gin/web/controller"
	"gin/web/middlewares"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var c *gin.Engine

func HttpRun() {
	c = gin.New()
	// 加载资源
	loadResource()

	// 路由注册
	registerRoute()

	// 启动http服务
	c.Run(util.Config.HttpAddr)
}

func registerRoute() {

	c.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, util.SuccessRespPackage("pong"))
		return
	})

	// 用户登陆
	c.POST("/admin/user/login", controller.LoginAdmin)
	c.GET("admin/test", controller.Test)

	admin := c.Group("/admin").Use(middlewares.AdminAuth())
	{
		// 创建用户
		admin.POST("/user/add", controller.AddAdminUser)
		admin.GET("/user/current", controller.Current)

		// 订单
		admin.GET("/order/list", controller.OrderList)
		admin.POST("/order/create", controller.AdminCreateOrder)
		admin.POST("/order/change/status", controller.ChangeOrderStatus)
		admin.POST("/order/edit/extra", controller.EditOrderExtra)
	}

	order := c.Group("/order").Use(middlewares.AdminAuth())
	{
		order.POST("/create/batch", web.BatchCreateOrder)
	}

	tools := c.Group("/tools").Use(middlewares.AdminAuth())
	{
		tools.POST("/file/upload", web.Upload)
	}

	img := c.Group("/img")
	{
		img.GET("/:imgName", web.ShowImage)
	}

	// 视图路由
	ResourceRegisterRoute(c)
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
