package db

import (
	"eebot/g"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlDB *gorm.DB

func InitMysql() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		g.Config.GetString("analysis.mysql.user"),
		g.Config.GetString("analysis.mysql.password"),
		g.Config.GetString("analysis.mysql.addr"),
		g.Config.GetString("analysis.mysql.databse"),
	)
	fmt.Printf("dsn: %v\n", dsn)
	SqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("panic: mysql, ", err.Error())
	}
	return
}
