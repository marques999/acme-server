package products

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

func List(database *gorm.DB) (int, interface{}) {
	products := []Product{}
	database.Find(&products)
	return http.StatusOK, products
}

func Insert(context *gin.Context, database *gorm.DB) (int, interface{}) {

	product := Product{}

	if ex := context.Bind(&product); ex != nil {
		return http.StatusBadRequest, common.JSON(ex)
	} else if ex := database.Save(&product).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusCreated, product
	}
}

func Find(context *gin.Context, database *gorm.DB) (int, interface{}) {

	product := Product{}

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if ex := database.First(&product, "barcode = ?", id).Error; ex != nil {
		return http.StatusNotFound, common.JSON(ex)
	} else {
		return http.StatusOK, product
	}
}
