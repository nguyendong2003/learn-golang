package storage

import (
	"context"
	"restfulapi/common"
	"restfulapi/module/item/model"
)

func (s *sqlStore) DeleteItem(ctx context.Context, condition map[string]interface{}) error {
	deletedStatus := model.ItemStatusDeleted

	if err := s.db.Table(model.TodoItem{}.TableName()).
		Where(condition).
		Updates(map[string]interface{}{"status": deletedStatus.String()}).Error; err != nil {

		return common.ErrDB(err)
	}

	return nil
}
