package entity

import (
	"errors"
	"main/common"
)

const (
	EntityName = "Item"
)

var (
	ErrTitleIsBlank = errors.New("Title cannot be blank")
	ErrItemDeleted = errors.New("Item has been deleted")
	ErrCannotCreateEntity = common.ErrCannotCreateEntity(EntityName, errors.New("cannot create entity"))
	ErrCannotUpdateEntity = common.ErrCannotUpdateEntity(EntityName, errors.New("cannot update entity"))
	ErrCannotDeleteEntity = common.ErrCannotDeleteEntity(EntityName, errors.New("cannot delete entity"))	
	ErrCannotGetEntity = common.ErrCannotGetEntity(EntityName, errors.New("cannot get entity"))
	ErrCannotListEntity = common.ErrCannotListEntity(EntityName, errors.New("cannot list entity"))
	ErrEntityDeleted = common.ErrEntityDeleted(EntityName, errors.New("entity has been deleted"))

)

type TodoItem struct {
	common.SQLModel
	Title       string      `json:"title" gorm:"column:title;"`
	Description string      `json:"description" gorm:"column:description;"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

type TodoItemCreation struct {
	Id          int         `json:"-" gorm:"column:id;"`
	Title       string      `json:"title" gorm:"column:title;"`
	Description string      `json:"description" gorm:"column:description;"`
	Status      *ItemStatus `json:"status" gorm:"column:status;"`
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}