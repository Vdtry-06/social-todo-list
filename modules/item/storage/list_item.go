package storage

import (
	"context"
	"main/common"
	"main/modules/item/entity"
)

func (sql *sqlStore) ListItem(
	ctx context.Context,
	filter *entity.Filter,
	paging *common.Paging,
	moreKeys ... string,
) ([]entity.TodoItem, error) {

	var results []entity.TodoItem

	// db := sql.db.Where("status <> ?", "Deleted")
	db := sql.db

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Table(entity.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}
	
	if err := db.Order("id desc").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&results).Error; err != nil {

		return nil, err
	}
	return results, nil
}