package customers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"net/http"
)

func List(context *gin.Context, database *gorm.DB) {
	customers := []Customer{}
	database.Preload("CreditCard").Find(&customers)
	context.JSON(http.StatusOK, customers)
}

func Authenticate(database *gorm.DB, username string, password string) (string, bool) {

	customer := Customer{}
	dbException := database.Preload("CreditCard").First(&customer, "username = ?", username).Error

	if dbException != nil {
		return username, false
	}

	if (username == "admin" && password == "admin") || (customer.Username == username && customer.Password == password) {
		return username, true
	}

	return username, false
}

func Insert(context *gin.Context, database *gorm.DB) {

	customer := Customer{}

	if jsonException := context.Bind(&customer); jsonException == nil {

		if customer.Username != "admin" {

			if dbException := database.Save(&customer).Error; dbException == nil {
				context.JSON(http.StatusCreated, customer)
			} else {
				common.RespondError(context, http.StatusInternalServerError, dbException.Error())
			}
		} else {
			common.RespondError(context, http.StatusUnauthorized, "permissionDenied")
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, jsonException.Error())
	}
}

func Find(context *gin.Context, database *gorm.DB, username string) {

	customerId, validParameters := context.Params.Get("id")

	if validParameters {

		if common.DenyAuthorization(context, username, customerId) {
			return
		}

		customer := GetCustomerOr404(context, database, customerId)

		if customer != nil {
			context.JSON(http.StatusOK, customer)
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func Update(context *gin.Context, database *gorm.DB, session string) {

	customerId, validParameters := context.Params.Get("id")

	if validParameters {

		if common.DenyAuthorization(context, session, customerId) {
			return
		}

		customer := GetCustomerOr404(context, database, customerId)

		if customer == nil {
			return
		}

		if jsonException := context.Bind(&customer); jsonException != nil {

			dbException := database.Preload("CreditCard").Update(&customer).Error

			if dbException == nil {
				context.JSON(http.StatusOK, customer)
			} else {
				common.RespondError(context, http.StatusInternalServerError, dbException.Error())
			}
		} else {
			common.RespondError(context, http.StatusBadRequest, jsonException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func Delete(context *gin.Context, database *gorm.DB, username string) {

	customerId, validParameters := context.Params.Get("id")

	if validParameters {

		if common.DenyAuthorization(context, username, customerId) {
			return
		}

		dbException := database.Delete(&Customer{Username: customerId}).Error

		if dbException == nil {
			context.JSON(http.StatusNoContent, nil)
		} else {
			common.RespondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		common.RespondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func GetCustomerOr404(context *gin.Context, database *gorm.DB, username string) *Customer {

	customer := Customer{}
	dbException := database.Preload("CreditCard").First(&customer, "username = ?", username).Error

	if dbException == nil {
		return &customer
	} else {
		common.RespondError(context, http.StatusNotFound, dbException.Error())
	}

	return nil
}
