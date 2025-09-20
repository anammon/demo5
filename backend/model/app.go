package model

import (
	"time"
)

type App struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Icon        string    `json:"icon" gorm:"type:varchar(255)"`   // 应用图标
	Type        string    `json:"type" gorm:"type:varchar(50)"`    // 应用类型
	Tags        string    `json:"tags" gorm:"type:varchar(255)"`   // 标签（逗号分隔）
	Status      string    `json:"status" gorm:"type:varchar(20)"`  // 状态（如 active/inactive）
	Author      string    `json:"author" gorm:"type:varchar(100)"` // 作者
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (App) TableName() string {
	return "appdb"
}
