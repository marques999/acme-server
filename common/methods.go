package common

import "github.com/gin-gonic/gin"

func JSON(ex error) map[string]interface{} {
	return gin.H{"error": ex.Error()}
}

func MissingParameter() map[string]interface{} {
	return gin.H{"error": "missingParameter"}
}

func PermissionDenied() map[string]interface{} {
	return gin.H{"error": "permissionDenied"}
}
