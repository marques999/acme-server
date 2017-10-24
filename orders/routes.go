package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func InitializeRoutes(database *gorm.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/orders")
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET("/", func(context *gin.Context) {
				context.JSON(List(database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.POST("/", func(context *gin.Context) {
				context.JSON(Insert(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.GET("/:id", func(context *gin.Context) {
				context.JSON(Find(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.PUT("/:id", func(context *gin.Context) {
				context.JSON(Checkout(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.DELETE("/:id", func(context *gin.Context) {
				context.JSON(Delete(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})
		}
	}
}