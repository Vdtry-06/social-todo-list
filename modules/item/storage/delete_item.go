package storage

import (
	"context"
	"main/common"
	"main/modules/item/entity"
)

func (sql *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {

	deletedStatus := entity.ItemStatusDeleted

	if err := sql.db.Table(entity.TodoItem{}.TableName()).
		Where(cond).
		Updates(map[string]interface{}{
		"status": deletedStatus.String(),
	}).Error; err != nil {
		return common.ErrDB(err)
	}
	
	return nil
}