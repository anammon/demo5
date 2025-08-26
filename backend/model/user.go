package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   //`gorm:"primaryKey"`
	Account   string `gorm:"unique"`
	Password  string
	Email     string `gorm:"unique"`
	Name      string `gorm:"unique"`
	CreatedAt time.Time
	Birthday  time.Time      `json:"birthday" time_format:"2006-01-02" time_utc:"0"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "userdb"
}
