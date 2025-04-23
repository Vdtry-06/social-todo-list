package main

import (
	"fmt"
	"log"
	"main/common"
	"main/middleware"
	"main/modules/item/transport/gin"
	"net/http"
	"os"

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

	r.Use(middleware.Recovery(db))

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
			items.GET("", ginItem.ListItem(db))
			items.GET("/:id", ginItem.GetItem(db))
			items.PATCH("/:id", ginItem.UpdateItem(db))
			items.DELETE("/:id", ginItem.DeleteItem(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		go func() {
			defer common.Recovery()
			fmt.Println([]int{}[0])
		}()
		
		c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		})
	})
	r.Run(":3000") // port 3000
}