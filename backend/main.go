package main

import (
	"backend/config"
	"backend/model"
	"backend/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	model.ConnectDatabase()
	r := gin.Default()
	err := model.DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("数据库迁移失败", err)
	}

	routers.UserRouters(r)
	r.Run()
}
