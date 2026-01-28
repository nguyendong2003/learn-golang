package business

import (
	"context"
	"restfulapi/common"
	"restfulapi/module/item/model"
)

type DeleteItemStorage interface {
	GetItem(ctx context.Context, condition map[string]interface{}) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, condition map[string]interface{}) error
}

type deleteItemBusiness struct {
	store DeleteItemStorage
}

func NewDeleteItemBusiness(store DeleteItemStorage) *deleteItemBusiness {
	return &deleteItemBusiness{store: store}
}

func (business *deleteItemBusiness) DeleteItemById(ctx context.Context, id int) error {
	data, err := business.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		if err == common.RecordNotFound {
			return common.ErrCannotGetEntity(model.EntityName, err)
		}

		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	if data.Status != nil && *data.Status == model.ItemStatusDeleted {
		return common.ErrEntityDeleted(model.EntityName, model.ErrItemDeleted)
	}

	if err := business.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}

	return nil
}
