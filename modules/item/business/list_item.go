package business

import (
	"context"
	"main/common"
	"main/modules/item/entity"
)

type ListItemStorage interface {
	ListItem(
		ctx context.Context, 
		filter *entity.Filter, 
		paging *common.Paging, 
		moreKeys ...string,
	) ([]entity.TodoItem, error)
}

type listItemBusiness struct {
	store ListItemStorage
}

func NewListItemBusiness(store ListItemStorage) *listItemBusiness {
	return &listItemBusiness{store: store}
}

func (business *listItemBusiness) ListItem(
	ctx context.Context, 
	filter *entity.Filter, 
	paging *common.Paging,
) ([]entity.TodoItem, error) {
	data, err := business.store.ListItem(ctx, filter, paging); 

	if err != nil {
		if err == common.RecordNotFound {
			return nil, common.ErrCannotListEntity(entity.EntityName, err)
		}
		return nil, err
	}
	return data, nil
}