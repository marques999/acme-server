package products

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/model"
	"net/http"
)

func List(context *gin.Context, database *gorm.DB) {
	products := []model.Product{}
	database.Find(&products)
	context.JSON(http.StatusOK, products)
}

func Insert(context *gin.Context, database *gorm.DB) {

	product := model.Product{}
	jsonException := context.Bind(&product)

	if jsonException == nil {

		dbException := database.Save(&product).Error

		if dbException == nil {
			context.JSON(http.StatusCreated, product)
		} else {
			respondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		respondError(context, http.StatusBadRequest, jsonException.Error())
	}
}

func Find(context *gin.Context, database *gorm.DB) {

	barcode, barcodeFound := context.Params.Get("id")

	if barcodeFound {

		product := GetProductOr404(context, database, barcode)

		if product != nil {
			context.JSON(http.StatusOK, product)
		}
	} else {
		respondError(context, http.StatusBadRequest, "invalidParameter")
	}
}

func Update(context *gin.Context, database *gorm.DB) {

	barcode, barcodeFound := context.Params.Get("id")

	if barcodeFound {

		product := GetProductOr404(context, database, barcode)

		if product == nil {
			return
		}

		jsonException := context.Bind(&product)

		if jsonException != nil {

			dbException := database.Update(&product).Error

			if dbException == nil {
				context.JSON(http.StatusOK, product)
			} else {
				respondError(context, http.StatusInternalServerError, dbException.Error())
			}
		} else {
			respondError(context, http.StatusBadRequest, jsonException.Error())
		}
	} else {
		respondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func Delete(context *gin.Context, database *gorm.DB) {

	barcode, barcodeFound := context.Params.Get("id")

	if barcodeFound {

		dbException := database.Delete(&model.Product{Barcode: barcode}).Error

		if dbException == nil {
			context.JSON(http.StatusNoContent, nil)
		} else {
			respondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		respondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func GetProductOr404(context *gin.Context, database *gorm.DB, barcode string) *model.Product {

	project := model.Product{}
	dbException := database.First(&project, model.Product{Barcode: barcode}).Error

	if dbException == nil {
		return &project
	} else {
		respondError(context, http.StatusNotFound, dbException.Error())
	}

	return nil
}

func respondError(context *gin.Context, statusCode int, message string) {
	context.JSON(statusCode, gin.H{"error": message})
}
