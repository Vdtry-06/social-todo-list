package storage

import (
	"context"
	"errors"
	"main/common"
	"main/modules/item/entity"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockStrategy interface {
	SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation)
	CheckResult(t *testing.T, err error)
}

// success
type SuccessStrategy struct{}

func (s *SuccessStrategy) SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `todo_items` (`title`,`description`,`status`) VALUES (?,?,?)")).
		WithArgs(data.Title, data.Description, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func (s *SuccessStrategy) CheckResult(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// error database
type DatabaseErrorStrategy struct{}

func (s *DatabaseErrorStrategy) SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation) {
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `todo_items`")).
		WillReturnError(common.ErrDB(errors.New("database error")))
	mock.ExpectRollback()
}

func (s *DatabaseErrorStrategy) CheckResult(t *testing.T, err error) {
	expectedErr := common.ErrDB(errors.New("database error"))
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got: %v", expectedErr, err)
	}
}



func TestCreateItem(t *testing.T) {

	type testCase struct {
		name          string
		data          *entity.TodoItemCreation
		strategy      MockStrategy
	}

	// Create database mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).
		WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.27"))

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{Conn: db}),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	store := NewSQLStore(gormDB)
	ctx := context.Background()

	testCases := []testCase{
		{
			name: "Success",
			data: &entity.TodoItemCreation{
				Title:       "Test Item",
				Description: "Test Description",
			},
			strategy: &SuccessStrategy{},
		},
		{
			name: "Error database",
			data: &entity.TodoItemCreation{
				Title:       "Test Item",
				Description: "Test Description",
			},
			strategy: &DatabaseErrorStrategy{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.strategy.SetupMock(mock, tc.data)

			err := store.CreateItem(ctx, tc.data)

			tc.strategy.CheckResult(t, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}