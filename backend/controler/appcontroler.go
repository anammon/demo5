package controler

import (
	"backend/config"
	"backend/model"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AppControler struct {
}

var cachekeylist = "applist"

func GetAppCacheKey(Id uint) string {
	return "app:" + strconv.Itoa(int(Id))
}
func (AppControler) Create(c *gin.Context) {
	var app model.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong param",
		})
		return
	}
	if strings.TrimSpace(app.Name) == "" || strings.TrimSpace(app.Description) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't be empty",
		})
		return
	}
	var exist model.App
	if err := model.DB.Where("name = ?", app.Name).First(&exist).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "the app name has existed",
		})
		return
	}
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	if err := model.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := config.RedisClient.Del(config.RedisCtx, cachekeylist).Err(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, app)
}
func (AppControler) Update(c *gin.Context) {
	var req model.App
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if req.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid app ID",
		})
		return
	}
	var app model.App
	if err := model.DB.First(&app, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 只允许更新部分字段
	app.Name = req.Name
	app.Description = req.Description
	app.Icon = req.Icon
	app.Type = req.Type
	app.Tags = req.Tags
	app.Status = req.Status
	app.Author = req.Author
	app.UpdatedAt = time.Now()
	if err := model.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := config.RedisClient.Del(config.RedisCtx, cachekeylist).Err(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	appkey := GetAppCacheKey(app.ID)
	if err := config.RedisClient.Del(config.RedisCtx, appkey).Err(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := config.RedisClient.Del(config.RedisCtx, cachekeylist).Err(); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, app)
}
func (AppControler) GetApp(c *gin.Context) {
	CachedData, err := config.RedisClient.Get(config.RedisCtx, cachekeylist).Result()
	var applist []model.App
	if err == redis.Nil {

		query := model.DB.Model(&model.App{}) //构建查询
		if name := c.Query("name"); name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		if description := c.Query("description"); description != "" {
			query = query.Where("description LIKE ?", "%"+description+"%")
		}
		if err := query.Find(&applist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		appjson, err := json.Marshal(applist)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := config.RedisClient.Set(config.RedisCtx, cachekeylist, appjson, 10*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		if err := json.Unmarshal([]byte(CachedData), &applist); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if name := c.Param("name"); name != "" {
			ChoseAftName := make([]model.App, 0)
			for _, appitem := range applist {
				if strings.Contains(strings.ToLower(appitem.Name), strings.ToLower(name)) {
					ChoseAftName = append(ChoseAftName, appitem)
				}
			}
			applist = ChoseAftName
		}
	}

	c.JSON(http.StatusOK, applist)
}
func (AppControler) GetAppById(c *gin.Context) {
	idstr := c.Param("app_id")
	id, err := strconv.Atoi(idstr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid app id",
		})
		return
	}
	appkey := GetAppCacheKey(uint(id))
	CachedData, err := config.RedisClient.Get(config.RedisCtx, appkey).Result()
	var app model.App
	if err == redis.Nil {
		if err := model.DB.First(&app, id).Error; err != nil { //redis.Nil--没找到缓存
			c.JSON(http.StatusNotFound, gin.H{
				"error": "app not found",
			})
			return
		}
		appjson, err := json.Marshal(app)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := config.RedisClient.Set(config.RedisCtx, appkey, appjson, 10*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if err != nil { //访问redis出错
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		if err := json.Unmarshal([]byte(CachedData), &app); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, app)
}
