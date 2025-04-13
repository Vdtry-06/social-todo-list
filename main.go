package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	fmt.Println("Hello, World!")

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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
		"message": item,
		})
	})
	r.Run(":3000") // port 3000
}