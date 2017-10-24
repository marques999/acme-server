package orders

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
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
	dbException := database.Preload("Products").First(&order, getQueryOptions(orderId, username)).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	order.Status = Purchased
	dbException = database.Update(&order).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	return http.StatusOK, generateJson(order)
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

	barcodeList := []string{}
	jsonException := json.Unmarshal(payload, &barcodeList)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	}

	customer := customers.Customer{}
	dbException := database.First(&customer, "username = ?", username).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	signature, _ := base64.StdEncoding.DecodeString(orderPOST.Signature)
	publicKey := decodePublicKey(customer.PublicKey)

	if publicKey == nil {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	cryptoException := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, sha1Checksum, signature)

	if cryptoException != nil {
		return http.StatusUnauthorized, common.JSON(cryptoException)
	}

	customerCart := []products.Product{}
	dbException = database.Where("barcode in (?)", barcodeList).Find(&customerCart).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	order := Order{Customer: customer.ID, Status: ValidationFailed, Total: CalculateTotal(customerCart)}
	dbException = database.Create(&order).Association("Products").Append(customerCart).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	orderToken, hashException := GenerateToken(&order)

	if hashException == nil {
		order.Token = orderToken
	} else {
		return http.StatusInternalServerError, common.JSON(hashException)
	}

	if customer.CreditCard.Validity.After(time.Now()) && rand.Float64() <= common.NaniProbability {
		order.Status = ValidationComplete
	}

	dbException = database.Save(&order).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	return http.StatusOK, generateJson(order)
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

	if dbException != nil {
		return http.StatusNotFound, common.JSON(dbException)
	}

	return http.StatusCreated, generateJson(order)
}

func Delete(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	orderId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	}

	dbException := database.Delete(&Order{}, getQueryOptions(orderId, username)).Error

	if dbException != nil {
		return http.StatusUnauthorized, common.JSON(dbException)
	}

	return http.StatusNoContent, nil
}