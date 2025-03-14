package initialize

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql(MysqlDataSource string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{})
	if err != nil {
		panic("mysql数据库连接失败, error=" + err.Error())
	} else {
		logx.Info("mysql数据库连接成功！")
	}
	return db
}
