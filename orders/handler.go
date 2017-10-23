package orders

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/products"
)

func List(database *gorm.DB, customerId string) (int, interface{}) {

	orders := []Order{}

	if customerId == common.AdminAccount {
		database.Preload("Products").Find(&orders)
	} else {
		database.Preload("Products").Find(&orders, "customer = ?", customerId)
	}

	return http.StatusOK, orders
}

func Insert(context *gin.Context, database *gorm.DB, session string) (int, interface{}) {

	jsonBody := common.Encrypted{}

	if ex := context.Bind(&jsonBody); ex != nil {
		return http.StatusBadRequest, common.JSON(ex)
	}

	sha1Checksum := encodeSha1([]byte(jsonBody.Payload))
	payload, base64Exception := base64.StdEncoding.DecodeString(jsonBody.Payload)

	if base64Exception != nil {
		return http.StatusBadRequest, common.JSON(base64Exception)
	}

	orderPOST := []string{}
	jsonException := json.Unmarshal(payload, &orderPOST)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	}

	customer := customers.Customer{}

	if ex := database.Preload("CreditCard").First(&customer, "username = ?", session).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	}

	signature, _ := base64.StdEncoding.DecodeString(jsonBody.Signature)
	publicKey := decodePublicKey(customer.PublicKey)

	if publicKey == nil {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	array := []products.Product{}
	order := Order{Customer: &customer, Valid: false}

	if ex := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, sha1Checksum, signature); ex != nil {
		return http.StatusUnauthorized, common.JSON(ex)
	} else if ex := database.Where("barcode in (?)", orderPOST).Find(&array).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := database.Create(&order).Association("Products").Append(array).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	}

	hashId, hashException := GenerateHashId(&order)

	if hashException != nil {
		return http.StatusInternalServerError, common.JSON(hashException)
	}

	order.Valid = true
	order.Token = hashId

	if tokenPayload, ex := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(hashId)); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := database.Save(&order).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, common.Encrypted{
			Signature: base64.StdEncoding.EncodeToString(encodeSha1([]byte(hashId))),
			Payload:   base64.StdEncoding.EncodeToString(tokenPayload),
		}
	}
}

func Find(context *gin.Context, database *gorm.DB, customerId string) (int, interface{}) {

	order := Order{}

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if ex := database.Preload("Products").First(&order, getQueryOptions(id, customerId)).Error; ex != nil {
		return http.StatusNotFound, common.JSON(ex)
	} else {
		return http.StatusCreated, order
	}
}

func Delete(context *gin.Context, database *gorm.DB, customerId string) (int, interface{}) {

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if ex := database.Delete(&Order{}, getQueryOptions(id, customerId)).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusNoContent, nil
	}
}
