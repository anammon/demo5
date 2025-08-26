package routers

import (
	"backend/controler"

	"github.com/gin-gonic/gin"
)

func UserRouters(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/register", controler.UserControler{}.Register)
		userGroup.POST("/login", controler.UserControler{}.Login)
	}

}
