package application

import (
	"time"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/customers"
)

func GetAuthenticator(database *gorm.DB) *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Unauthorized: func(context *gin.Context, statusCode int, message string) {
			context.JSON(statusCode, gin.H{"error": message})
		},
		Authenticator: func(username string, password string, context *gin.Context) (string, bool) {
			return customers.Authenticate(database, username, password)
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}
