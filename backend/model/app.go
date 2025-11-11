package model

import (
	"time"
)

type App struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
    Icon        string    `json:"icon" gorm:"type:varchar(255)"`
    Type        string    `json:"type" gorm:"type:varchar(50)"`
	Tags        string    `json:"tags" gorm:"type:varchar(255)"`   
	Status      string    `json:"status" gorm:"type:varchar(20)"`  
	Author      string    `json:"author" gorm:"type:varchar(100)"` 
	Likes       int       `json:"likes" gorm:"type:int;default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (App) TableName() string {
	return "appdb"
}
