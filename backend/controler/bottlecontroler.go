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
	"gorm.io/gorm"
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
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	expireAt := tomorrow.Unix()
	luaScript := `
			local key = KEYS[1]
			local expire_at = tonumber(ARGV[1])
			local limit = 3
			
			local current = redis.call("GET", key)
			if current then
				current = tonumber(current)
				if current >= limit then
					return {0, 0}  -- 已超限
				end
			end
			
			local new_count = redis.call("INCR", key)
			if new_count == 1 then
				redis.call("EXPIREAT", key, expire_at)
			end
			
			return {1, limit - new_count}  -- 成功
		`
	result, err := config.RedisClient.Eval(config.RedisCtx, luaScript, []string{key}, expireAt).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		return
	}

	results, _ := result.([]interface{})
	success, _ := results[0].(int64)
	remaining, _ := results[1].(int64)

	if success == 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error":     "今日捡瓶子次数已用完，明天再来吧~",
			"remaining": remaining,
		})
		return
	}

	// 在数据库事务中使用悲观锁捡瓶子
	var pickedBottle model.Bottle
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		// 使用FOR UPDATE锁定记录，防止并发捡到同一个瓶子
		var availableBottles []model.Bottle
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("is_picked = ?", false).
			Limit(10). // 只取10个即可
			Find(&availableBottles).Error; err != nil {
			return err
		}

		if len(availableBottles) == 0 {
			return gorm.ErrRecordNotFound
		}
		b := availableBottles[rand.Intn(len(availableBottles))]
		if err := tx.Model(&b).Updates(map[string]interface{}{
			"is_picked":    true,
			"pick_user_id": pickUID,
		}).Error; err != nil {
			return err
		}

		pickedBottle = b
		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			config.RedisClient.Decr(config.RedisCtx, key)
			DefaultBottle := []string{
				"今天没有捡到漂流瓶~ 给你一句：加油！你很棒 💪",
				"学习很累、生活很难，但你也很强 🌈",
				"没瓶子？那就给你一个拥抱 🤗",
			}
			msg := DefaultBottle[rand.Intn(len(DefaultBottle))]
			c.JSON(http.StatusOK, gin.H{
				"id":         0,
				"content":    msg,
				"created_at": time.Now(),
				"is_system":  true,
				"remaining":  remaining,
			})
			return
		}
		config.RedisClient.Decr(config.RedisCtx, key)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"id":         pickedBottle.ID,
		"content":    pickedBottle.Content,
		"created_at": pickedBottle.CreatedAt,
		"is_picked":  pickedBottle.IsPicked,
		"remaining":  remaining,
	}

	if !pickedBottle.IsAnonymous {
		response["is_anonymous"] = false
		response["throw_user_info"] = pickedBottle.ThrowUserID
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误,应为2000-01-01"})
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
