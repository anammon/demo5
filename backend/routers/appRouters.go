package routers

import (
	"backend/controler"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AppRouters(r *gin.Engine) {
	appGroup := r.Group("/app").Use(middlewares.AuthMiddleware())
	{
		appGroup.POST("/create", controler.AppControler{}.Create)
		appGroup.POST("/getapp", controler.AppControler{}.Create)
	}
}
