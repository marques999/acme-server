package orders

import (
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func InitializeRoutes(database *gorm.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/orders")
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET("/", func(context *gin.Context) {
				context.JSON(List(database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})

			routes.POST("/", func(context *gin.Context) {
				context.JSON(Insert(context, database))
			})

			routes.GET("/:id", func(context *gin.Context) {
				context.JSON(Find(context, database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})

			routes.PUT("/:id", func(context *gin.Context) {
				context.JSON(Validate(context, database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})

			routes.DELETE("/:id", func(context *gin.Context) {
				context.JSON(Delete(context, database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})
		}
	}
}
