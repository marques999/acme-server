package customers

import (
	"fmt"
	"net/http"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

func InitializeRoutes(database *gorm.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	routes := router.Group("/customers")
	{
		routes.Use(middleware.MiddlewareFunc())
		{
			routes.GET("/", func(context *gin.Context) {

				if jwt.ExtractClaims(context)["id"] == common.AdminAccount {
					context.JSON(List(database))
				} else {
					context.JSON(http.StatusUnauthorized, nil)
				}
			})

			routes.GET("/:id", func(context *gin.Context) {
				context.JSON(Find(context, database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})

			routes.POST("/", func(context *gin.Context) {
				context.JSON(Insert(context, database))
			})

			routes.PUT("/:id", func(context *gin.Context) {
				context.JSON(Update(context, database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})

			routes.DELETE("/:id", func(context *gin.Context) {
				context.JSON(Delete(context, database, fmt.Sprint(jwt.ExtractClaims(context)["id"])))
			})
		}
	}
}
