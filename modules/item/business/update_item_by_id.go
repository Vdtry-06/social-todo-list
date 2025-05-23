package business

import (
	"context"
	"main/common"
	"main/modules/item/entity"
)

type UpdateItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *entity.TodoItemUpdate) error
}

type updateItemBusiness struct {
	store UpdateItemStorage
}

func NewUpdateItemBusiness(store UpdateItemStorage) *updateItemBusiness {
	return &updateItemBusiness{store: store}
}

func (business *updateItemBusiness) UpdateItemById(ctx context.Context, id int, dataUpdate *entity.TodoItemUpdate) error {
	data, err := business.store.GetItem(ctx, map[string]interface{}{"id": id}); 

	if err != nil {
		if err == common.RecordNotFound {
			return common.ErrCannotUpdateEntity(entity.EntityName, err)
		}
		return common.ErrCannotUpdateEntity(entity.EntityName, err)
	}

	if data.Status != nil && *data.Status == entity.ItemStatusDeleted {	
		return common.ErrEntityDeleted(entity.EntityName, entity.ErrItemDeleted)
	}

	if err := business.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return common.ErrCannotUpdateEntity(entity.EntityName, err)
	}

	return nil
}