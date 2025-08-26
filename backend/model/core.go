package model

import (
	"backend/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := config.AppConfig.Database.Dsn
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
