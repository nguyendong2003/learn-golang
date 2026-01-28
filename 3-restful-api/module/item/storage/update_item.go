package storage

import (
	"context"
	"restfulapi/common"
	"restfulapi/module/item/model"
)

func (s *sqlStore) UpdateItem(ctx context.Context, condition map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := s.db.Where(condition).Updates(dataUpdate).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
