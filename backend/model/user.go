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
	Birthday  time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "userdb"
}
