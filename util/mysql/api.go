package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

func InitMysql(dbUser, dbPassword, dbHost, dbName string, dbPort, maxIdleConns, maxOpenConns, connMaxLifetime int) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level: Silent、Error、Warn、Info
			Colorful:      false,         // 禁用彩色打印
		},
	)

	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
	})
	if err != nil {
		log.Fatalln("open mysql failure!", dsn)
	}

	//开启连接池
	sqlDB, err := Db.DB()
	sqlDB.SetMaxIdleConns(maxIdleConns)                      //最大空闲连接
	sqlDB.SetMaxOpenConns(maxOpenConns)                      //最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime)) //最大生存时间(s)
}
