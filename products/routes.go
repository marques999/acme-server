package products

import (
	"net/http"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
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

			routes.POST("/", func(context *gin.Context) {

				if jwt.ExtractClaims(context)["id"] == common.AdminAccount {
					context.JSON(Insert(context, database))
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})
		}
	}
}
