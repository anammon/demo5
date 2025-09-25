package main

import (
	"backend/config"
	"backend/model"
	"backend/routers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()
	model.ConnectDatabase()
	r := gin.Default()
	r.Use(cors.Default())
	if err := model.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("数据库迁移失败: ", err)
	}
	routers.UserRouters(r)
	routers.AppRouters(r)
	if err := r.Run(config.AppConfig.App.Port); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}
