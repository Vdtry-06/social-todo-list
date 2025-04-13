package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

type TodoItem struct {
	Id 			int 		`json:"id"`
	Title 		string 		`json:"title"`
	Description string 		`json:"description"`
	Status 		string 		`json:"status"`
	CreatedAt 	*time.Time 	`json:"created_at"`
	UpdatedAt 	*time.Time 	`json:"updated_at,omitempty"`

}

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
	
	fmt.Println(db)

	now := time.Now().UTC()

	item := TodoItem{
		Id: 1,
		Title: "Learn Go",
		Description: "Learn Go programming language",
		Status: "In Progress",
		CreatedAt: &now,
		UpdatedAt: &now,
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
			items.POST("")
			items.GET("")
			items.GET("/:id")
			items.PATCH("/:id")
			items.DELETE("/:id")
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
		"message": item,
		})
	})
	r.Run(":3000") // port 3000
}