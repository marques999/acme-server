package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

func InitializeRoutes(middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	router.POST("/login", middleware.LoginHandler)

	routes := router.Group("/auth").Use(middleware.MiddlewareFunc())
	{
		routes.GET("/refresh_token", middleware.RefreshHandler)
	}
}