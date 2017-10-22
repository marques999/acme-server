package admin

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(database *gorm.DB, router *gin.Engine) {

	routes := router.Group("/admin")
	{
		routes.POST("/populate", func(context *gin.Context) {
			context.JSON(Populate(database))
		})

		routes.POST("/clean", func(context *gin.Context) {
			context.JSON(Clean(database))
		})
	}
}
