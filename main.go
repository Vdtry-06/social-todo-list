package main

import (
	"encoding/json"
	"fmt"
	"time"
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
		UpdatedAt: nil,
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonData))

	jsonStr := "{\"id\":1,\"title\":\"Learn Go\",\"description\":\"Learn Go programming language\",\"status\":\"In Progress\",\"created_at\":\"2023-10-01T12:00:00Z\",\"updated_at\":\"2023-10-01T12:00:00Z\"}"

	var item2 TodoItem

	if err := json.Unmarshal([]byte(jsonStr), &item2); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item2)
}