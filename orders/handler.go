package orders

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/products"
)

func List(database *gorm.DB, username string) (int, interface{}) {

	orders := []Order{}

	if username == common.AdminAccount {
		database.Preload("Products").Find(&orders)
	} else {
		database.Preload("Products").Find(&orders, "customer = ?", username)
	}

	jsonOrders := make([]map[string]interface{}, len(orders))

	for i, order := range orders {
		jsonOrders[i] = generateJson(order)
	}

	return http.StatusOK, jsonOrders
}

func Checkout(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	orderId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	}

	order := Order{}
	dbsException := database.Preload("Products").First(&order, getQueryOptions(orderId, username)).Error

	if dbsException != nil {
		return http.StatusInternalServerError, common.JSON(dbsException)
	}

	order.Status = 2
	dbuException := database.Update(&order).Error

	if dbuException == nil {
		return http.StatusOK, generateJson(order)
	} else {
		return http.StatusInternalServerError, common.JSON(dbuException)
	}
}

func Insert(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	orderPOST := OrderPOST{}
	bindException := context.Bind(&orderPOST)

	if bindException != nil {
		return http.StatusBadRequest, common.JSON(bindException)
	}

	sha1Checksum := encodeSha1([]byte(orderPOST.Payload))
	payload, base64Exception := base64.StdEncoding.DecodeString(orderPOST.Payload)

	if base64Exception != nil {
		return http.StatusBadRequest, common.JSON(base64Exception)
	}

	orderProducts := []string{}
	jsonException := json.Unmarshal(payload, &orderProducts)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	}

	customer := customers.Customer{}
	dbsException := database.First(&customer, "username = ?", username).Error

	if dbsException != nil {
		return http.StatusInternalServerError, common.JSON(dbsException)
	}

	signature, _ := base64.StdEncoding.DecodeString(orderPOST.Signature)
	publicKey := decodePublicKey(customer.PublicKey)

	if publicKey == nil {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	array := []products.Product{}
	order := Order{Customer: customer.ID, Status: 0}

	if ex := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, sha1Checksum, signature); ex != nil {
		return http.StatusUnauthorized, common.JSON(ex)
	} else if ex := database.Where("barcode in (?)", orderProducts).Find(&array).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := database.Create(&order).Association("Products").Append(array).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	}

	if token, ex := GenerateToken(&order); ex == nil {
		order.Status = 1
		order.Token = token
	} else {
		return http.StatusInternalServerError, common.JSON(ex)
	}

	if ex := database.Save(&order).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, generateJson(order)
	}
}

func Find(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	if username != common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	tokenId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	}

	order := Order{}
	dbException := database.Preload("Products").First(&order, "token = ?", tokenId).Error

	if dbException == nil {
		return http.StatusCreated, generateJson(order)
	} else {
		return http.StatusNotFound, common.JSON(dbException)
	}
}

func Delete(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	orderId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	}

	dbException := database.Delete(&Order{}, getQueryOptions(orderId, username)).Error

	if dbException == nil {
		return http.StatusNoContent, nil
	} else {
		return http.StatusUnauthorized, common.JSON(dbException)
	}
}