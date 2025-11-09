package model

import "time"

// Bottle 漂流瓶模型
type Bottle struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Content     string    `json:"content" gorm:"type:varchar(2000);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	ThrowUserID uint32    `json:"throw_user_id" gorm:"not null;index"`
	PickUserID  *uint32   `json:"pick_user_id" gorm:"index"` // 可为空
	IsPicked    bool      `json:"is_picked" gorm:"default:false"`
	IsAnonymous bool      `json:"is_anonymous"`
}

func (Bottle) TableName() string {
	return "bottledb"
}
