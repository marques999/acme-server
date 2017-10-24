package products

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func InitializeRoutes(database *gorm.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/products")
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET("/", func(context *gin.Context) {
				context.JSON(List(database))
			})

			routes.GET("/:id", func(context *gin.Context) {
				context.JSON(Find(context, database))
			})
		}
	}
}