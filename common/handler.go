package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondError(context *gin.Context, statusCode int, message string) {
	context.JSON(statusCode, gin.H{"error": message})
}

func DenyAuthorization(context *gin.Context, username string, customerId string) bool {

	if username != "admin" && username != customerId {
		RespondError(context, http.StatusUnauthorized, "permissionDenied")
	} else {
		return false
	}

	return true
}
