package middleware

import (
	"main/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Recovery(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		
		defer func() {
			if err := recover(); err != nil {
				if err, ok := err.(error); ok {
					c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrInternal(err))
				}
				panic(err) // log the error and continue the program
			}
		}()

		c.Next()
	}
}