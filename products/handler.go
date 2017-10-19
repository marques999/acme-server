package products

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"net/http"
)

func List(context *gin.Context, database *gorm.DB) {
	products := []Product{}
	database.Find(&products)
	context.JSON(http.StatusOK, products)
}

func ListByBarcode(database *gorm.DB, barcode []string) []Product {
	products := []Product{}
	database.Where("barcode in (?)", barcode).Find(&products)
	return products
}

func Insert(context *gin.Context, database *gorm.DB) {

	product := Product{}
	jsonException := context.Bind(&product)

	if jsonException == nil {

		dbException := database.Save(&product).Error

		if dbException == nil {
			context.JSON(http.StatusCreated, product)
		} else {
			common.RespondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, jsonException.Error())
	}
}

func Find(context *gin.Context, database *gorm.DB) {

	productId, validParameters := context.Params.Get("id")

	if validParameters {

		product := Product{}
		dbException := database.First(&product, "barcode = ?", productId).Error

		if dbException == nil {
			context.JSON(http.StatusOK, product)
		} else {
			common.RespondError(context, http.StatusNotFound, dbException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, "invalidParameter")
	}
}
