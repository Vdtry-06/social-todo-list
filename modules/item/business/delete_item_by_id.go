package business

import (
	"context"
	"main/common"
	"main/modules/item/entity"
)

type DeleteItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type deleteItemBusiness struct {
	store DeleteItemStorage
}

func NewDeleteItemBusiness(store DeleteItemStorage) *deleteItemBusiness {
	return &deleteItemBusiness{store: store}
}

func (business *deleteItemBusiness) DeleteItemById(ctx context.Context, id int) error {
	data, err := business.store.GetItem(ctx, map[string]interface{}{"id": id}); 

	if err != nil {
		if err == common.RecordNotFound {
			return common.ErrCannotDeleteEntity(entity.EntityName, err)
		}
		return err
	}

	if data.Status != nil && *data.Status == entity.ItemStatusDeleted {	
		return entity.ErrItemDeleted
	}

	if err := business.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return err
	}

	return nil
}