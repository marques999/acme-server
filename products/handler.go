package products

import (
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
)

func List(database *sqlx.DB) (int, interface{}) {

	if products, errors := getProducts(database); errors != nil {
		return http.StatusInternalServerError, errors
	} else {
		return http.StatusOK, products
	}
}

func Find(context *gin.Context, database *sqlx.DB) (int, interface{}) {

	if barcode, exists := context.Params.Get(common.Id); exists == false {
		return common.MissingParameter()
	} else if product, errors := getProduct(database, barcode); errors != nil {
		return http.StatusNotFound, common.JSON(errors)
	} else {
		return http.StatusOK, product
	}
}

func Insert(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	productJson := ProductInsert{}

	if username != common.AdminAccount {
		return common.PermisssionDenied()
	} else if errors := context.Bind(&productJson); errors != nil {
		return http.StatusBadRequest, common.JSON(errors)
	} else if product, errors := insertProduct(database, productJson); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, product.GenerateJson()
	}
}

func Update(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	productJson := ProductUpdate{}

	if username != common.AdminAccount {
		return common.PermisssionDenied()
	} else if barcode, exists := context.Params.Get(common.Id); exists == false {
		return common.MissingParameter()
	} else if errors := context.Bind(&productJson); errors != nil {
		return http.StatusBadRequest, common.JSON(errors)
	} else if product, errors := updateProduct(database, barcode, productJson); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, product.GenerateJson()
	}
}

func Delete(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	if username != common.AdminAccount {
		return common.PermisssionDenied()
	} else if barcode, exists := context.Params.Get(common.Id); exists == false {
		return common.MissingParameter()
	} else if _, errors := deleteProduct(database, barcode); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusNoContent, nil
	}
}