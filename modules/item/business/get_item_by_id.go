package business

import (
	"context"
	"main/common"
	"main/modules/item/entity"
)

type GetItemStorage interface {
	GetItem(ctx context.Context, cond map[string] interface{}) (*entity.TodoItem, error)
}

type getItemBusiness struct {
	store GetItemStorage
}

func NewGetItemBusiness(store GetItemStorage) *getItemBusiness {
	return &getItemBusiness{store: store}
}

func (business *getItemBusiness) GetItemById(ctx context.Context, id int) (*entity.TodoItem, error) {
	data, err := business.store.GetItem(ctx, map[string]interface{}{"id": id}); 

	if err != nil {
		return nil, common.ErrCannotGetEntity(entity.EntityName, err)
	}
	return data, nil
}