package main

import (
	"log"
	"main/common"
	"main/modules/item/entity"
	ginItem "main/modules/item/transport/gin"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
  	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err) // log.Fatal will print the error and exit the program
	}

	r := gin.Default()

	// CRUD: Create, Read, Update, Delete
	// POST: /v1/items (create a new item)
	// GET: /v1/items (list items) /v1/items/items?page=1 || /v1/items?cursor=fdsfsdk
	// GET: /v1/items/:id (get item detail by id)
	// PUT || PATCH: /v1/items/:id (update a item by id)
	// DELETE: /v1/items/:id (delete a item by id)

	v1 := r.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", ginItem.CreateItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id", UpdateItem(db))
			items.DELETE("/:id", DeleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		})
	})
	r.Run(":3000") // port 3000
}


func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data entity.TodoItem

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data entity.TodoItemUpdate

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := db.Table(entity.TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
			"status": "Deleted",
		}).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging
		
		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		paging.Process()

		var results []entity.TodoItem

		db = db.Where("status <> ?", "Deleted")

		if err := db.Table(entity.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		
		if err := db.Order("id desc").
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, common.NewSuccessResponse(results, paging, nil))
	}
}