package business

import (
	"context"
	"restfulapi/module/item/model"
)

type UpdateItemStorage interface {
	GetItem(ctx context.Context, condition map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, condition map[string]interface{}, dataUpdate *model.TodoItemUpdate) error
}

type updateItemBusiness struct {
	store UpdateItemStorage
}

func NewUpdateItemBusiness(store UpdateItemStorage) *updateItemBusiness {
	return &updateItemBusiness{store: store}
}

func (business *updateItemBusiness) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := business.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return err
	}

	if data.Status != nil && *data.Status == model.ItemStatusDeleted {
		return model.ErrItemDeleted
	}

	if err := business.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return err
	}

	return nil
}
