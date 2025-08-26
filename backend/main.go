package main

import (
	"backend/model"
	"backend/routers"
	"log"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	err:=model.DB.AutoMigrate(&model.User{})
	if err!=nil{
		log.Fatal("数据库迁移失败",err)
	}
	
	routers.UserRouters(r)
 
}
