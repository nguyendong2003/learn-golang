package business

import (
	"context"
	"restfulapi/common"
	"restfulapi/module/item/model"
)

type ListItemStorage interface {
	ListItem(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type listItemBusiness struct {
	store ListItemStorage
}

func NewListItemBusiness(store ListItemStorage) *listItemBusiness {
	return &listItemBusiness{store: store}
}

func (business *listItemBusiness) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.TodoItem, error) {
	data, err := business.store.ListItem(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.EntityName, err)
	}

	return data, nil
}
