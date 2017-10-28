package orders

import (
	"net/http"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
)

func List(database *sqlx.DB, username string) (int, interface{}) {

	if orders, errors := listOrders(database, username); errors != nil {
		return http.StatusInternalServerError, errors
	} else {
		return http.StatusOK, orders
	}
}

func Find(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	if token, exists := context.Params.Get("id"); exists == false {
		return common.MissingParameter()
	} else if order, errors := getOrder(database, getQueryOptions(token, username)); errors != nil {
		return http.StatusNotFound, common.JSON(errors)
	} else if customerCart, errors := getCustomerCart(database, order.ID); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusCreated, order.generateJson(customerCart)
	}
}

func Insert(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	orderPOST := OrderPOST{}

	if errors := context.Bind(&orderPOST); errors != nil {
		return http.StatusBadRequest, common.JSON(errors)
	} else if jsonProducts, errors := json.Marshal(orderPOST.Products); errors != nil {
		return http.StatusBadRequest, common.JSON(errors)
	} else if customer, errors := customers.GetCustomer(database, username); errors != nil {
		return http.StatusUnauthorized, common.JSON(errors)
	} else if errors = verifySignature(
		customer.PublicKey, orderPOST.Signature, encodeSha1(jsonProducts),
	); errors != nil {
		return http.StatusUnauthorized, errors
	} else if order, errors := insertOrder(
		database, customer, orderPOST.Products...,
	); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, order
	}
}

func Purchase(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	if token, exists := context.Params.Get("id"); exists == false {
		return common.MissingParameter()
	} else if order, ex := updateOrder(database, getQueryOptions(token, username), Purchased); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if customerCart, ex := getCustomerCart(database, order.ID); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, order.generateJson(customerCart)
	}
}

func Delete(context *gin.Context, database *sqlx.DB, customer string) (int, interface{}) {

	if token, exists := context.Params.Get("id"); exists == false {
		return common.MissingParameter()
	} else if _, ex := deleteOrder(database, getQueryOptions(token, customer)); ex != nil {
		return http.StatusUnauthorized, common.JSON(ex)
	} else {
		return http.StatusNoContent, nil
	}
}