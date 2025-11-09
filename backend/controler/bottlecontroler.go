package controler

import (
	"backend/config"
	"backend/model"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func toUint32(v interface{}) (uint32, error) {
	switch t := v.(type) {
	case uint32:
		return t, nil
	case uint:
		return uint32(t), nil
	case uint64:
		return uint32(t), nil
	case int:
		if t < 0 {
			return 0, fmt.Errorf("negative int")
		}
		return uint32(t), nil
	case int32:
		if t < 0 {
			return 0, fmt.Errorf("negative int32")
		}
		return uint32(t), nil
	case int64:
		if t < 0 {
			return 0, fmt.Errorf("negative int64")
		}
		return uint32(t), nil
	case float64:
		if t < 0 {
			return 0, fmt.Errorf("negative float64")
		}
		return uint32(t), nil
	case string:
		if t == "" {
			return 0, fmt.Errorf("empty string")
		}
		u, err := strconv.ParseUint(t, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(u), nil
	default:
		return 0, fmt.Errorf("unsupported userID type %T", v)
	}
}

type BottleController struct {
}

type ThrowRequest struct {
	Content     string `json:"content" binding:"required"`
	IsAnonymous bool   `json:"is_anonymous"`
}

// 扔瓶子
func (BottleController) ThrowBottle(c *gin.Context) {
	var rep ThrowRequest
	if err := c.ShouldBind(&rep); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if len([]rune(rep.Content)) > 600 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too long, 最多600个字!"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	throwUID, err := toUint32(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	bottle := model.Bottle{
		Content:     rep.Content,
		ThrowUserID: throwUID,
		IsAnonymous: rep.IsAnonymous,
		IsPicked:    false,
		PickUserID:  nil, // 初始未被捡
	}

	if err := model.DB.Create(&bottle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bottle)
}

// 捡瓶子
func (BottleController) PickBottle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	pickUID, err := toUint32(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	key := fmt.Sprintf("pick_limit:%d:%s", pickUID, time.Now().Format("2006-01-02"))
	count, err := config.RedisClient.Get(config.RedisCtx, key).Int()
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		return
	}
	if count >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "今日捡瓶子次数已用完，明天再来吧~", "remaining": 0})
		return
	}

	// 只查未被捡的
	var availableBottles []model.Bottle
	if err := model.DB.Where("is_picked = ?", false).Find(&availableBottles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(availableBottles) == 0 {
		DefaultBottle := []string{
			"今天没有捡到漂流瓶~ 给你一句：加油！你很棒 💪",
			"生活很难，但你也很强 🌈",
			"没瓶子？那就给你一个拥抱 🤗",
		}
		msg := DefaultBottle[rand.Intn(len(DefaultBottle))]
		c.JSON(http.StatusOK, gin.H{"id": 0, "content": msg, "created_at": time.Now(), "is_system": true})
		return
	}

	// 随机取一个
	b := availableBottles[rand.Intn(len(availableBottles))]
	b.IsPicked = true
	tmp := pickUID
	b.PickUserID = &tmp

	if err := model.DB.Save(&b).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := config.RedisClient.Incr(config.RedisCtx, key).Result(); err == nil && count == 0 {
		tomorrow := time.Now().Add(24 * time.Hour)
		midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
		config.RedisClient.ExpireAt(config.RedisCtx, key, midnight)
	}

	// 响应
	response := gin.H{
		"id":         b.ID,
		"content":    b.Content,
		"created_at": b.CreatedAt,
		"is_picked":  b.IsPicked,
	}
	if !b.IsAnonymous {
		response["is_anonymous"] = false
		response["throw_user_info"] = b.ThrowUserID
	} else {
		response["is_anonymous"] = true
	}

	c.JSON(http.StatusOK, response)
}

// 我的扔瓶历史
func (BottleController) GetMyThrownBottles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	throwUID, err := toUint32(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	query := model.DB.Where("throw_user_id = ?", throwUID)
	if datestr := c.Query("date"); datestr != "" {
		if rightdate, err := time.Parse("2006-01-02", datestr); err == nil {
			query = query.Where("created_at between ? and ?", rightdate, rightdate.Add(24*time.Hour))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误,应为YYYY-MM-DD"})
			return
		}
	}

	var bottles []model.Bottle
	if err := query.Order("created_at DESC").Find(&bottles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bottles)
}

// 我的捡瓶历史
func (BottleController) GetMyPickedBottles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	pickUID, err := toUint32(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	query := model.DB.Where("pick_user_id = ?", pickUID)
	if datestr := c.Query("date"); datestr != "" {
		if rightdate, err := time.Parse("2006-01-02", datestr); err == nil {
			query = query.Where("created_at between ? and ?", rightdate, rightdate.Add(24*time.Hour))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误,应为YYYY-MM-DD"})
			return
		}
	}

	var bottles []model.Bottle
	if err := query.Order("created_at DESC").Find(&bottles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bottles)
}
