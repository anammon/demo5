package service

import (
	"backend/config"
	"backend/model"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var likeQueue = make(chan string, 10000)

func EnqueueLike(appID string) {
	select {
	case likeQueue <- appID:
	default:
		log.Println("like queue full, dropping like for app:", appID)
	}
}

func StartLikeFlusher(db *gorm.DB, interval time.Duration) {
	// 初始化确保Redis中至少存在每个app的likes 键
	var apps []model.App
	if err := db.Find(&apps).Error; err == nil {
		for _, app := range apps {
			key := "app:" + strconv.FormatUint(uint64(app.ID), 10) + ":likes"
			exists, err := config.RedisClient.Exists(config.RedisCtx, key).Result()
			if err != nil {
				continue
			}
			if exists == 0 {
				config.RedisClient.Set(config.RedisCtx, key, strconv.Itoa(app.Likes), 0)
			}
		}
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		pending := make(map[string]struct{})
		for {
			select {
			case id := <-likeQueue:
				pending[id] = struct{}{}
			case <-ticker.C:
				if len(pending) == 0 {
					continue
				}
				ids := make([]string, 0, len(pending))
				for k := range pending {
					ids = append(ids, k)
				}
				for _, id := range ids {
					key := "app:" + id + ":likes"
					s, err := config.RedisClient.Get(config.RedisCtx, key).Result()
					if err != nil {
						if err.Error() == "redis: nil" {
							continue
						}
						log.Println("error reading redis key for likes", key, err)
						continue
					}
					val, err := strconv.Atoi(s)
					if err != nil {
						log.Println("invalid like value in redis for key", key, s)
						continue
					}
					// 写到数据库
					err = db.Model(&model.App{}).Where("id = ?", id).Update("likes", val).Error
					if err != nil {
						log.Println("failed to update likes for app", id, err)
						continue
					}
				}
				pending = make(map[string]struct{}) //清空pending
			}
		}
	}()
}
