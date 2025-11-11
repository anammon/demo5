package controler

import (
	"backend/config"
	"backend/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AppControler struct {
}

// 单条应用缓存 key
func GetAppCacheKey(Id uint) string {
	return "app:" + strconv.Itoa(int(Id))
}
const appListCachePrefix = "apps:page:"  // 前缀
const appListCacheTTL = 10 * time.Minute // 缓存过期时间

func buildAppListCacheKey(page, pageSize int, name, description string) string {
	ne := url.QueryEscape(name)
	de := url.QueryEscape(description)
	return fmt.Sprintf("%s%d:size:%d:name:%s:desc:%s", appListCachePrefix, page, pageSize, ne, de)
}

// 使用 SCAN + DEL，避免 KEYS 在大数据量下阻塞
func deleteAppListCache() error {
	pattern := appListCachePrefix + "*" 
	var cursor uint64 = 0
	for {
		keys, cur, err := config.RedisClient.Scan(config.RedisCtx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := config.RedisClient.Del(config.RedisCtx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = cur
		if cursor == 0 {
			break
		}
	}
	return nil
}

// 分页响应结构
type AppListPage struct {
	Data     []model.App `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
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

	if err := deleteAppListCache(); err != nil {
		log.Printf("deleteAppListCache error after Create: %v", err)
	}

	// 同时删除单条缓存
	if err := config.RedisClient.Del(config.RedisCtx, GetAppCacheKey(app.ID)).Err(); err != nil {
		log.Printf("delete single app cache after Create error: %v", err)
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
	if err := deleteAppListCache(); err != nil {
		log.Printf("deleteAppListCache error after Update: %v", err)
	}
	appkey := GetAppCacheKey(app.ID)
	if err := config.RedisClient.Del(config.RedisCtx, appkey).Err(); err != nil {
		log.Printf("delete single app cache after Update error: %v", err)
	}

	c.JSON(http.StatusOK, app)
}
func (AppControler) GetApp(c *gin.Context) {
	name := c.Query("name")
	description := c.Query("description")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}

	cacheKey := buildAppListCacheKey(page, pageSize, name, description)
	if cached, err := config.RedisClient.Get(config.RedisCtx, cacheKey).Result(); err == nil {
		var pageRes AppListPage
		if err := json.Unmarshal([]byte(cached), &pageRes); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"data":     pageRes.Data,
				"total":    pageRes.Total,
				"page":     pageRes.Page,
				"pageSize": pageRes.PageSize,
			})
			return
		}
	} else if err != redis.Nil {
		log.Printf("redis get error for key %s: %v", cacheKey, err)
	}
    
	// 2) 缓存未命中从 DB 查询
	var applist []model.App
	query := model.DB.Model(&model.App{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&applist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3) 构建响应并缓存
	pageRes := AppListPage{
		Data:     applist,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	if bs, err := json.Marshal(pageRes); err == nil {
		if err := config.RedisClient.Set(config.RedisCtx, cacheKey, bs, appListCacheTTL).Err(); err != nil {
			log.Printf("redis set error for key %s: %v", cacheKey, err)
		}
	} else {
		log.Printf("json.Marshal pageRes error: %v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"data":     applist,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
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
			// 记录缓存写入错误但不影响正常响应
			log.Printf("redis set single app cache error for key %s: %v", appkey, err)
		}
	} else if err != nil { //访问redis出错
		log.Printf("redis get single app key %s error: %v", appkey, err)
		// 继续从 DB 读取
		if err := model.DB.First(&app, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "app not found",
			})
			return
		}
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
