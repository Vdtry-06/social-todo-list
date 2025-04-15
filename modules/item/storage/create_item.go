package storage

import (
	"context"
	"main/modules/item/entity"
)

func (sql *sqlStore) CreateItem(ctx context.Context, data *entity.TodoItemCreation) error {
	if err := sql.db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}