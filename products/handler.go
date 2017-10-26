package products

import (
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
)

func LIST(database *sqlx.DB) (int, interface{}) {

	products, sqlException := getProducts(database)

	if sqlException != nil {
		return http.StatusInternalServerError, sqlException
	}

	productsJson := make([]ProductJSON, len(products))

	for index, product := range products {
		productsJson[index] = product.GenerateJson()
	}

	return http.StatusOK, productsJson
}

func GET(context *gin.Context, database *sqlx.DB) (int, interface{}) {

	barcode, paramExists := context.Params.Get(common.Id)

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	}

	product, dbException := getProduct(database, barcode)

	if dbException != nil {
		return http.StatusNotFound, common.JSON(dbException)
	}

	return http.StatusOK, product.GenerateJson()
}

func Insert(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	productJson := ProductJSON{}
	jsonException := context.Bind(&productJson)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	} else if username != common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	product, sqlException := insertProduct(database, productJson)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	} else {
		return http.StatusOK, product.GenerateJson()
	}
}

func Update(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	productJson := ProductJSON{}
	jsonException := context.Bind(&productJson)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	} else if username != common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	product, sqlException := updateProduct(database, productJson)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	} else {
		return http.StatusOK, product.GenerateJson()
	}
}

func Delete(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	barcode, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	_, dbException := deleteProduct(database, barcode)

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	} else {
		return http.StatusNoContent, nil
	}
}