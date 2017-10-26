package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func InitializeRoutes(database *sqlx.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/orders")
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET("/", func(context *gin.Context) {
				context.JSON(LIST(database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.POST("/", func(context *gin.Context) {
				context.JSON(INSERT(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.GET("/:id", func(context *gin.Context) {
				context.JSON(GET(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.PUT("/:id", func(context *gin.Context) {
				context.JSON(PUT(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})

			routes.DELETE("/:id", func(context *gin.Context) {
				context.JSON(DELETE(context, database, (jwt.ExtractClaims(context)["id"]).(string)))
			})
		}
	}
}