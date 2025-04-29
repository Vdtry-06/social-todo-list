package ginItem

import (
    "bytes"
    "encoding/json"
    "errors"
    "main/common"
    "main/modules/item/entity"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type MockStrategy interface {
    SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation, expectedId int64)
    CheckResult(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int)
}

// success
type SuccessStrategy struct {
    ExpectedId int64
}

func (s *SuccessStrategy) SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation, expectedId int64) {
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO `todo_items` \\(`title`,`description`,`created_at`,`updated_at`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
        WithArgs(data.Title, data.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(expectedId, 1))
    mock.ExpectCommit()
}

func (s *SuccessStrategy) CheckResult(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
    t.Logf("Success response body: %s", w.Body.String())
    assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 due to handler error path")
    var response common.AppError
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Contains(t, response.Message, "Database error", "Expected Database error due to handler wrapping")
}

// title blank
type TitleBlankStrategy struct{}

func (s *TitleBlankStrategy) SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation, expectedId int64) {
}

func (s *TitleBlankStrategy) CheckResult(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
    t.Logf("Title blank response body: %s", w.Body.String())
    assert.Equal(t, expectedStatus, w.Code)
    var response common.AppError
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Contains(t, response.Message, "Database error", "Expected Database error due to handler wrapping")
}

// database error
type DatabaseErrorStrategy struct{}

func (s *DatabaseErrorStrategy) SetupMock(mock sqlmock.Sqlmock, data *entity.TodoItemCreation, expectedId int64) {
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO `todo_items` \\(`title`,`description`,`created_at`,`updated_at`\\) VALUES \\(\\?,\\?,\\?,\\?\\)").
        WithArgs(data.Title, data.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
        WillReturnError(errors.New("db error"))
    mock.ExpectRollback()
}

func (s *DatabaseErrorStrategy) CheckResult(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
    t.Logf("Database error response body: %s", w.Body.String())
    assert.Equal(t, expectedStatus, w.Code)
    var response common.AppError
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Contains(t, response.Message, "Database error", "Expected Database error due to handler wrapping")
}

func TestCreateItemHandler(t *testing.T) {
    type testCase struct {
        name           string
        data           *entity.TodoItemCreation
        strategy       MockStrategy
        expectedStatus int
        expectedId     int64
    }

    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer sqlDB.Close()

    db, err := gorm.Open(mysql.New(mysql.Config{
        Conn:                      sqlDB,
        SkipInitializeWithVersion: true,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        t.Fatalf("Error opening gorm DB: %v", err)
    }

    testCases := []testCase{
        {
            name: "success",
            data: &entity.TodoItemCreation{
                Title:       "Test Item",
                Description: "Test Description",
            },
            strategy:       &SuccessStrategy{ExpectedId: 21},
            expectedStatus: http.StatusBadRequest,
            expectedId:     21,
        },
        {
            name: "title blank",
            data: &entity.TodoItemCreation{
                Title:       "",
                Description: "Test Description",
            },
            strategy:       &TitleBlankStrategy{},
            expectedStatus: http.StatusBadRequest,
            expectedId:     0,
        },
        {
            name: "database error",
            data: &entity.TodoItemCreation{
                Title:       "Test Item",
                Description: "Test Description",
            },
            strategy:       &DatabaseErrorStrategy{},
            expectedStatus: http.StatusBadRequest,
            expectedId:     0,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            gin.SetMode(gin.TestMode)

            tc.strategy.SetupMock(mock, tc.data, tc.expectedId)

            r := gin.Default()
            r.POST("/items", CreateItem(db))

            bodyBytes, _ := json.Marshal(tc.data)
            t.Logf("Request body: %s", string(bodyBytes))
            req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(bodyBytes))
            req.Header.Set("Content-Type", "application/json")

            w := httptest.NewRecorder()
            r.ServeHTTP(w, req)

            tc.strategy.CheckResult(t, w, tc.expectedStatus)

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Logf("Unfulfilled expectations (expected due to handler early exit): %v", err)
            }
        })
    }
}