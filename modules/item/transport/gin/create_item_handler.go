package ginItem

import (
	"main/common"
	"main/modules/item/business"
	"main/modules/item/entity"
	"main/modules/item/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data entity.TodoItemCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		store:= storage.NewSQLStore(db)
		biz := business.NewCreateItemBusiness(store)

		if err := biz.CreateNewItem(c.Request.Context(), &data); err != nil {
			if err := c.ShouldBind(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}