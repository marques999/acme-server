package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func InitializeRoutes(database *sqlx.DB, middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	router.POST("/login", middleware.LoginHandler)

	routes := router.Group("/auth").Use(middleware.MiddlewareFunc())
	{
		routes.GET("/refresh_token", middleware.RefreshHandler)

		router.POST("/reset", func(context *gin.Context) {
			context.JSON(Clean(database))
		})
	}
}