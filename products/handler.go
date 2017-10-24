package products

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
)

func List(database *gorm.DB) (int, interface{}) {

	products := []Product{}
	database.Find(&products)
	jsonProducts := make([]map[string]interface{}, len(products))

	for i, product := range products {

		jsonProducts[i] = gin.H{
			"name":    product.Name,
			"brand":   product.Brand,
			"price":   product.Price,
			"barcode": product.Barcode,
			"uri":     product.ImageUri,
		}
	}

	return http.StatusOK, jsonProducts
}

func Find(context *gin.Context, database *gorm.DB) (int, interface{}) {

	barcode, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	}

	product := Product{}
	dbException := database.First(&product, "barcode = ?", barcode).Error

	if dbException != nil {
		return http.StatusNotFound, common.JSON(dbException)
	}

	return http.StatusOK, gin.H{
		"name":        product.Name,
		"brand":       product.Brand,
		"price":       product.Price,
		"barcode":     product.Barcode,
		"uri":         product.ImageUri,
		"description": product.Description,
	}
}