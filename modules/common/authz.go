package common

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")

func EnsureGroupMember(db *gorm.DB, groupID, userID uuid.UUID) error {
	var count int64
	if err := db.Table("group_members").
		Where("group_id = ? AND user_id = ? AND is_active = true", groupID, userID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrForbidden
	}
	return nil
}
