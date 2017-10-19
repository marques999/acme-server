package tokens

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"net/http"
	"github.com/marques999/acme-server/products"
)

func List(context *gin.Context, database *gorm.DB, customerId string) {

	tokens := []Token{}

	if customerId == "admin" {
		database.Preload("Products").Find(&tokens)
	} else {
		database.Preload("Products").Find(&tokens, "customer = ?", customerId)
	}

	context.JSON(http.StatusOK, tokens)
}

func Insert(context *gin.Context, database *gorm.DB) {

	jsonBody := TokenPOST{}
	jsonException := context.Bind(&jsonBody)

	if jsonException == nil {

		token := Token{Customer: jsonBody.Customer}

		dbException := database.Create(&token).Association("Products").Append(
			products.ListByBarcode(database, jsonBody.Products),
		).Error

		if dbException == nil {
			context.JSON(http.StatusCreated, token)
		} else {
			common.RespondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, jsonException.Error())
	}
}

func Find(context *gin.Context, database *gorm.DB, customerId string) {

	tokenId, validParameters := context.Params.Get("id")

	if validParameters {

		token := Token{}
		dbException := database.Preload("Products").First(&token, getQueryOptions(tokenId, customerId)).Error

		if dbException == nil {
			context.JSON(http.StatusOK, token)
		} else {
			common.RespondError(context, http.StatusNotFound, dbException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func Delete(context *gin.Context, database *gorm.DB, customerId string) {

	tokenId, validParameters := context.Params.Get("id")

	if validParameters {

		dbException := database.Delete(&Token{}, getQueryOptions(tokenId, customerId)).Error

		if dbException == nil {
			context.JSON(http.StatusNoContent, nil)
		} else {
			common.RespondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func getQueryOptions(tokenId string, customerId string) map[string]interface{} {
	if customerId == "admin" {
		return map[string]interface{}{
			"id": tokenId,
		}
	} else {
		return map[string]interface{}{
			"id":       tokenId,
			"customer": customerId,
		}
	}
}
