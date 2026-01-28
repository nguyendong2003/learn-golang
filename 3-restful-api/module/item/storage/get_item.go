package storage

import (
	"context"
	"restfulapi/common"
	"restfulapi/module/item/model"

	"gorm.io/gorm"
)

func (s *sqlStore) GetItem(ctx context.Context, condition map[string]interface{}) (*model.TodoItem, error) {
	var data model.TodoItem

	if err := s.db.Where(condition).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &data, nil
}
