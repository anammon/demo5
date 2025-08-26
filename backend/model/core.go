package model

import (
	"backend/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/userdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	sqldb, err2 := db.DB()
	sqldb.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConns)
	sqldb.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConns)
	sqldb.SetConnMaxLifetime(time.Hour)
	DB = db
	if err2 != nil {
		panic("failed to config database")
	}
    
}
