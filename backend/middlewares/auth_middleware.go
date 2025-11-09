package middlewares

import (
	"backend/model"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供token",
			})
			c.Abort() //没有权限，终止后续处理
			return
		}
		account, err := utils.PaserJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalide token",
			})
			c.Abort() //没有权限，终止后续处理
			return
		}
		var user model.User
		if err := model.DB.Where("account = ?", account).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "用户不存在",
			})
			c.Abort()
			return
		}
		c.Set("account", account)
		c.Set("userID", user.ID)
		c.Next()
	}

}
