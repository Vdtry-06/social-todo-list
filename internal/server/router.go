package server

import (
	"fmt"
	"main/common"
	ginItem "main/modules/item/transport/gin"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
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
}