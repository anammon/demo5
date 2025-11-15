package service

import (
	"errors"

	"gorm.io/gorm"
)

type RelationService struct {
	db *gorm.DB
}

func NewRelationService(db *gorm.DB) *RelationService {
	return &RelationService{
		db: db,
	}
}

// follow,unfollow,getfollowlist
func (rs *RelationService) Follow(FollowerId, FollowingID uint) error {
	if FollowerId == FollowingID {
		return errors.New("不能关注自己")
	}

}
