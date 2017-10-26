package common

import (
	"github.com/gin-gonic/gin"
	"github.com/Masterminds/squirrel"
)

func JSON(ex error) map[string]interface{} {
	return gin.H{"error": ex.Error()}
}

func MissingParameter() map[string]interface{} {
	return gin.H{"error": "missingParameter"}
}

func PermissionDenied() map[string]interface{} {
	return gin.H{"error": "permissionDenied"}
}

func HasPermissions(username string, customerId string) bool {
	return username == AdminAccount || username != customerId
}

func StatementBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}