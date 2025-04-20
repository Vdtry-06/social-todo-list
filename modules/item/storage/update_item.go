package storage

import (
	"context"
	"main/modules/item/entity"
)

func (sql *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *entity.TodoItemUpdate) error {


	if err := sql.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return err
	}

	return nil
}