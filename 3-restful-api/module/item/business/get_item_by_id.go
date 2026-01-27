package business

import (
	"context"
	"restfulapi/module/item/model"
)

type GetItemStorage interface {
	GetItem(ctx context.Context, condition map[string]interface{}) (*model.TodoItem, error)
}

type getItemBusiness struct {
	store GetItemStorage
}

func NewGetItemBusiness(store GetItemStorage) *getItemBusiness {
	return &getItemBusiness{store: store}
}

func (business *getItemBusiness) GetItemById(ctx context.Context, id int) (*model.TodoItem, error) {
	data, err := business.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return nil, err
	}

	return data, nil
}
