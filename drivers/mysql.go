package drivers

import (
	"context"
	"fmt"
	"gin/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)


const (
	MysqlPoolMaxIdeNum  = 10  // mysql 连接池最大空闲数
	MysqlPoolMaxOpenNum = 100 // mysql 连接池最大连接数
)

var client *gorm.DB


func Mysql() *gorm.DB {
	return client.WithContext(context.Background())
}

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local",
		util.Config.MysqlUser,
		util.Config.MysqlPassword,
		util.Config.MysqlHOST,
		util.Config.MysqlPort,
		util.Config.MysqlDatabase,
	)
	var __logger logger.Interface
	if util.Config.Debug {
		__logger = logger.Default.LogMode(logger.Info)
	} else {
		__logger = logger.Default
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      __logger,
		PrepareStmt: true,
		//SkipDefaultTransaction: true, // 默认事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		util.Log.Fatal("failed to connection mysql server : %v", err)
	}

	originDB, err := db.DB()
	if err != nil {
		util.Log.Fatal("failed to init mysql driver : %v", err)
	}
	originDB.SetMaxIdleConns(MysqlPoolMaxIdeNum)
	originDB.SetMaxOpenConns(MysqlPoolMaxOpenNum)
	client = db

}