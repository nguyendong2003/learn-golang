package storage

import (
	"context"
	"restfulapi/common"
	"restfulapi/module/item/model"
)

// tham số đầu tiên nên là context.Context
func (s *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
