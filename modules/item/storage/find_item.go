package storage

import (
	"context"
	"main/modules/item/entity"
)

func (sql *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*entity.TodoItem, error) {
	
	var data entity.TodoItem

	if err := sql.db.Where(cond).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}