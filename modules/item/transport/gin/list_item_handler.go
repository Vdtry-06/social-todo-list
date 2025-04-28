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

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging
		customLimitPage := 5
		
		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
		
		paging.Process(customLimitPage)

		var filter entity.Filter
		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))
			return
		}
	
		store := storage.NewSQLStore(db)
		biz := business.NewListItemBusiness(store)

		results, err := biz.ListItem(c.Request.Context(), &filter, &paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(results, paging, filter))
	}
}