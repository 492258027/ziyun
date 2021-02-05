package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

const (
	dbUser     string = "root"
	dbPassword string = "123456"
	dbHost     string = "192.168.73.3"
	dbPort     int    = 3306
	dbName     string = "mytest"
)

var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	dbUser, dbPassword, dbHost, dbPort, dbName)

type User struct {
	gorm.Model
	//Id       int `gorm:"primary_key" json:"id"`
	Userid   string
	Username string
	Password string
	//Authorities []string
}

func main() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level: Silent、Error、Warn、Info
			Colorful:      false,         // 禁用彩色打印
		},
	)
	//db, err := gorm.Open("mysql", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
	})

	if err != nil {
		log.Fatalln("open mysql failure!")
	}

	//开启连接池
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)    //最大空闲连接
	sqlDB.SetMaxOpenConns(100)   //最大连接数
	sqlDB.SetConnMaxLifetime(30) //最大生存时间(s)

	//创建表
	db.AutoMigrate(&User{})
	// 将 "ENGINE=InnoDB" 添加到创建 `User` 的 SQL 里去
	//db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

	user := new(User)

	//添加数据
	user.Userid = "1001"
	user.Username = "xuheng"
	user.Password = "1234"
	//user.Authorities = []string{"read", "write"}
	if err := db.Create(user).Error; err != nil {
		log.Println("create failure")
	} else {
		log.Println("create success, ID:", user.Model.ID)
	}

	//修改用户
	if err := db.Model(user).Update("Username", "xuziyun").Error; err != nil {
		log.Println("modify failure")
	}

	//删除数据(软删除)
	if err := db.Delete(user).Error; err != nil {
		log.Println("logic delete failure")
	} else {
		log.Println("logic delete success")
	}
	//删除数据(真删除)
	user.ID = 1
	if err := db.Unscoped().Delete(user).Error; err != nil {
		log.Println("delete failure")
	} else {
		log.Println("delete success")
	}

	//查找数据
	db.Where("username = ? AND userid = ?", "xuheng", "1001").Find(&user)
	log.Println(user)

	var user_s []User
	db.Where("username IN (?)", []string{"xuheng", "xuziyun"}).Find(&user_s)
	for _, v := range user_s {
		log.Println(v)
	}

	//等同于 First
	db.Select("username").Where("id > 2").Order("username").Limit(1).Find(&user)
	log.Println(user)

	//join

	//事务
}
