package controler

import (
	"backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type LikeControler struct {
}

func (LikeControler) LikeApp(c *gin.Context) {
	appid := c.Param("app_id")
	likekey := "app:" + appid + ":likes"
	if err := config.RedisClient.Incr(config.RedisCtx, likekey).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "like success",
	})
}

func (LikeControler) GetAppLikes(c *gin.Context) {
	appid := c.Param("app_id")
	likekey := "app:" + appid + ":likes"
	likes, err := config.RedisClient.Get(config.RedisCtx, likekey).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, gin.H{
				"app_id": appid,
				"like":   0,
			}) //虽然err不为空，但其实操作正确，只是redis没这个键，返回0like就行
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}) //连接错误等等(real mistake)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"app_id": appid,
		"likes":  likes,
	})
}
