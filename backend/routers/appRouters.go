package routers

import (
	"backend/controler"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AppRouters(r *gin.Engine) {
	appGroup := r.Group("/app").Use(middlewares.AuthMiddleware())
	{
		// 应用相关
		appGroup.POST("/", controler.AppControler{}.Create)
		appGroup.GET("/", controler.AppControler{}.GetApp)
		appGroup.PUT("/:app_id", controler.AppControler{}.Update)
		appGroup.GET("/:app_id", controler.AppControler{}.GetAppById)
		appGroup.POST("/:app_id/like", controler.LikeControler{}.LikeApp)
		appGroup.GET("/:app_id/likes", controler.LikeControler{}.GetAppLikes)

		// 漂流瓶相关
		appGroup.POST("/bottle", controler.BottleController{}.ThrowBottle)                 // 扔瓶子
		appGroup.GET("/bottle/pick", controler.BottleController{}.PickBottle)              // 捡瓶子
		appGroup.GET("/bottle/my/thrown", controler.BottleController{}.GetMyThrownBottles) // 我的扔瓶子
		appGroup.GET("/bottle/my/picked", controler.BottleController{}.GetMyPickedBottles) // 我的捡瓶子
		// 矩阵计算相关
		appGroup.POST("/matrix/addition", controler.HandleMatrixAddition)
		appGroup.POST("/matrix/subtraction", controler.HandleMatrixSubtraction)
		appGroup.POST("/matrix/multiplication", controler.HandleMatrixMultiplication)
	}
}
