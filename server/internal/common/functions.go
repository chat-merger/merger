package common

import (
	"gorm.io/gorm"

	"github.com/chat-merger/merger/server/internal/model"
)

func Applications(db *gorm.DB, excludeIDs ...int) ([]*model.Application, error) {
	var apps []*model.Application
	if err := db.
		Where("id NOT IN (?)", excludeIDs).
		Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}
