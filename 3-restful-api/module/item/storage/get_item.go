package storage

import (
	"context"
	"restfulapi/module/item/model"
)

func (s *sqlStore) GetItem(ctx context.Context, condition map[string]interface{}) (*model.TodoItem, error) {
	var data model.TodoItem

	if err := s.db.Where(condition).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}
