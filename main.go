/*
 * @Date: 2022-08-18 11:02:36
 */
package main

import (
	"gin/drivers"
	"gin/util"
	"gin/web/routes"
)

func main() {
	// 初始化日志
	util.InitLog()
	// 初始化配置文件
	util.InitConfig()
	// 初始化数据库
	drivers.InitMysql()
	// http服务

	routes.HttpRun()
}
