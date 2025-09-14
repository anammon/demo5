package routers

import (
	"backend/controler"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AppRouters(r *gin.Engine) {
	appGroup := r.Group("/app").Use(middlewares.AuthMiddleware())
	{
		appGroup.POST("/", controler.AppControler{}.Create)
		appGroup.GET("/", controler.AppControler{}.GetApp)
		appGroup.POST("/:app_id/like", controler.LikeControler{}.LikeApp)
		appGroup.GET("/:app_id/likes", controler.LikeControler{}.GetAppLikes)
	}
}
