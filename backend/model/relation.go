package model

import "time"

type UserRelation struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	FollowerID  uint      `json:"follower_id" gorm:"not null;index"`  // 粉丝
	FollowingID uint      `json:"following_id" gorm:"not null;index"` // 博主
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserRelation) TableName() string {
	return "user_relation"
}
