package auth

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(middleware *jwt.GinJWTMiddleware, router *gin.Engine) {

	router.POST("/login", middleware.LoginHandler)

	routes := router.Group("/auth").Use(middleware.MiddlewareFunc())
	{
		routes.GET("/refresh_token", middleware.RefreshHandler)
	}
}
