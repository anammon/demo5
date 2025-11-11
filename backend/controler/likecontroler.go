package controler

import (
	"backend/config"
	"backend/model"
	"backend/service"
	"net/http"
	"strconv"

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
	service.EnqueueLike(appid)
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
			// redis 中没有这个键，尝试落库读取持久化值
			var app model.App
			if dbErr := model.DB.Where("id = ?", appid).First(&app).Error; dbErr == nil {
				c.JSON(http.StatusOK, gin.H{
					"app_id": appid,
					"likes":  app.Likes,
				})
				return
			}
			// 如果数据库也没有，就返回0
			c.JSON(http.StatusOK, gin.H{
				"app_id": appid,
				"likes":  0,
			}) //虽然err不为空，但其实操作正确，只是redis没这个键，返回0like就行
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}) //连接错误等等(real mistake)
		return
	}
	// 成功从 redis 读取，转换为数字返回
	if n, err := strconv.Atoi(likes); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"app_id": appid,
			"likes":  n,
		})
		return
	}
	// 若转换失败，返回字符串形式作为兜底
	c.JSON(http.StatusOK, gin.H{
		"app_id": appid,
		"likes":  likes,
	})
}
