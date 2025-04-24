package business

import (
	"context"
	"errors"
	"main/common"
	"main/modules/item/entity"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockCreateItemStorage struct {
	mock.Mock
}

func (m *MockCreateItemStorage) CreateItem(ctx context.Context, data *entity.TodoItemCreation) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type MockStrategy interface {
	SetupMock(mock *MockCreateItemStorage, data *entity.TodoItemCreation, ctx context.Context)
	CheckResult(t *testing.T, err error)
}

// success
type SuccessStrategy struct{}

func (s *SuccessStrategy) SetupMock(mock *MockCreateItemStorage, data *entity.TodoItemCreation, ctx context.Context) {
	mock.On("CreateItem", ctx, data).Return(nil).Once()
}

func (s *SuccessStrategy) CheckResult(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// title blank
type TitleBlankStrategy struct{}

func (s *TitleBlankStrategy) SetupMock(mock *MockCreateItemStorage, data *entity.TodoItemCreation, ctx context.Context) {
}

func (s *TitleBlankStrategy) CheckResult(t *testing.T, err error) {
	if err != entity.ErrTitleIsBlank {
		t.Errorf("Expected error %v, got: %v", entity.ErrTitleIsBlank, err)
	}
}

// database error
type DatabaseErrorStrategy struct{}

func (s *DatabaseErrorStrategy) SetupMock(mock *MockCreateItemStorage, data *entity.TodoItemCreation, ctx context.Context) {
	mock.On("CreateItem", ctx, data).Return(common.ErrDB(errors.New("db error"))).Once()
}

func (s *DatabaseErrorStrategy) CheckResult(t *testing.T, err error) {
	expectedErr := common.ErrCannotCreateEntity(entity.EntityName, errors.New("db error"))
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error %v, got: %v", expectedErr, err)
	}
}

func TestCreateNewItem(t *testing.T) {
	type testCase struct {
		name     string
		data     *entity.TodoItemCreation
		strategy MockStrategy
	}

	mockStore := &MockCreateItemStorage{}
	biz := NewCreateItemBusiness(mockStore)
	ctx := context.Background()

	testCases := []testCase{
		{
			name: "Success",
			data: &entity.TodoItemCreation{
				Title:       "Test item",
				Description: "Test description",
			},
			strategy: &SuccessStrategy{},
		},
		{
			name: "Title blank",
			data: &entity.TodoItemCreation{
				Title:       "",
				Description: "Test Description",
			},
			strategy: &TitleBlankStrategy{},
		},
		{
			name: "Database error",
			data: &entity.TodoItemCreation{
				Title:       "Test Item",
				Description: "Test Description",
			},
			strategy: &DatabaseErrorStrategy{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.strategy.SetupMock(mockStore, tc.data, ctx)

			err := biz.CreateNewItem(ctx, tc.data)

			tc.strategy.CheckResult(t, err)

			mockStore.AssertExpectations(t)
		})
	}
}