package main

import (
	"backend/config"
	"backend/middlewares"
	"backend/model"
	"backend/routers"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()
	model.ConnectDatabase()
	r := gin.Default()
	r.Use(middlewares.CorsMiddleware())
	// 数据库自动迁移
	if err := model.DB.AutoMigrate(&model.User{}, &model.App{}, &model.Bottle{}); err != nil {
		log.Fatal("数据库迁移失败: ", err)
	}
	routers.UserRouters(r)
	routers.AppRouters(r)
	r.Static("/assets", "D:/gin/demo5/frontend/dist/assets")
	r.Static("/icons", "D:/gin/demo5/frontend/dist/icons")
	r.StaticFile("/favicon.ico", "D:/gin/demo5/frontend/dist/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.File("D:/gin/demo5/frontend/dist/index.html")
	})
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/user") || strings.HasPrefix(path, "/app") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		c.File("D:/gin/demo5/frontend/dist/index.html")
	})
	port := config.AppConfig.App.Port
	var listenAddr string
	if strings.HasPrefix(port, ":") {
		listenAddr = "0.0.0.0" + port
	} else {
		listenAddr = "0.0.0.0:" + port
	}
	log.Printf("服务启动: %s", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}
