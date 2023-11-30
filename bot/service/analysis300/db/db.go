package db

import (
	"eebot/g"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SqlDB *gorm.DB

func HasPartition() bool {
	return SqlDB.Migrator().HasTable(PlayerPartition{})
}

func InitMysql() (err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 5, // Slow SQL threshold
			LogLevel:                  logger.Warn,     // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  false,           // Disable color
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		g.Config.GetString("analysis.mysql.user"),
		g.Config.GetString("analysis.mysql.password"),
		g.Config.GetString("analysis.mysql.addr"),
		g.Config.GetString("analysis.mysql.database"),
	)
	SqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic("panic: mysql, ", err.Error())
	}
	return
}
