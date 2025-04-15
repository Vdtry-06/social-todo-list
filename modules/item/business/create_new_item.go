package business

import (
	"context"
	"main/modules/item/entity"
	"strings"
)

type CreateItemStorage interface {
	CreateItem(ctx context.Context, data *entity.TodoItemCreation) error
}

type createItemBusiness struct {
	store CreateItemStorage
}

func NewCreateItemBusiness(store CreateItemStorage) *createItemBusiness {
	return &createItemBusiness{store: store}
}

func (business *createItemBusiness) CreateNewItem(ctx context.Context, data *entity.TodoItemCreation) error {
	title := strings.TrimSpace(data.Title)
	if title == "" {
		return entity.ErrTitleIsBlank
	}

	if err := business.store.CreateItem(ctx, data); err != nil {
		return err
	}
	return nil
}