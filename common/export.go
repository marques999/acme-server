package common

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/appleboy/gin-jwt"
	"github.com/Masterminds/squirrel"
)

var MissingParameterError = errors.New("request has missing or invalid \"id\" parameter")
var PermissionDeniedError = errors.New("user does not have permission to access the requested resource")

func GeneratePassword(original string) (string, error) {

	if hashed, exception := bcrypt.GenerateFromPassword(
		[]byte(original), bcrypt.DefaultCost,
	); exception == nil {
		return string(hashed), nil
	} else {
		return "", exception
	}
}

func JSON(exception error) map[string]interface{} {
	return gin.H{"error": exception.Error()}
}

func MissingParameter() (int, interface{}) {
	return http.StatusBadRequest, JSON(MissingParameterError)
}

func PermisssionDenied() (int, interface{}) {
	return http.StatusUnauthorized, JSON(PermissionDeniedError)
}

func HasPermissions(username string, customerId string) bool {
	return username == AdminAccount || username != customerId
}

func ParseId(context *gin.Context) string {
	return jwt.ExtractClaims(context)[Id].(string)
}

func SqlBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}