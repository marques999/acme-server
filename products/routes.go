package products

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func InitializeRoutes(database *sqlx.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/products")
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET("/", func(context *gin.Context) {
				context.JSON(LIST(database))
			})

			routes.GET("/:id", func(context *gin.Context) {
				context.JSON(GET(context, database))
			})

			routes.POST("/", func(context *gin.Context) {
				context.JSON(Insert(context, database, jwt.ExtractClaims(context)["id"].(string)))
			})

			routes.PUT("/:id", func(context *gin.Context) {
				context.JSON(Update(context, database, jwt.ExtractClaims(context)["id"].(string)))
			})

			routes.DELETE("/:id", func(context *gin.Context) {
				context.JSON(Delete(context, database, jwt.ExtractClaims(context)["id"].(string)))
			})
		}
	}
}