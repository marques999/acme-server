package customers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-backend/model"
	"net/http"
)

func List(context *gin.Context, database *gorm.DB, session string) {

	if session == "admin" {
		customers := []model.Customer{}
		database.Preload("CreditCard").Find(&customers)
		context.JSON(http.StatusOK, customers)
	} else {
		respondError(context, http.StatusUnauthorized, "permissionDenied")
	}
}

func Insert(context *gin.Context, database *gorm.DB) {

	customer := model.Customer{}
	jsonException := context.Bind(&customer)

	if jsonException == nil {

		dbException := database.Preload("CreditCard").Save(&customer).Error

		if dbException == nil {
			context.JSON(http.StatusCreated, customer)
		} else {
			respondError(context, http.StatusInternalServerError, dbException.Error())
		}
	} else {
		respondError(context, http.StatusBadRequest, jsonException.Error())
	}
}

func Find(context *gin.Context, database *gorm.DB, session string) {

	username, usernameFound := context.Params.Get("id")

	if usernameFound {

		if session == "admin" || session == username {

			customer := GetCustomerOr404(context, database, username)

			if customer != nil {
				context.JSON(http.StatusOK, customer)
			}
		} else {
			respondError(context, http.StatusUnauthorized, "permissionDenied")
		}
	}
}

func Update(context *gin.Context, database *gorm.DB, session string) {

	username, usernameFound := context.Params.Get("id")

	if usernameFound {

		if session == "admin" || session == username {

			customer := GetCustomerOr404(context, database, username)

			if customer == nil {
				return
			}

			jsonException := context.Bind(&customer)

			if jsonException != nil {

				dbException := database.Preload("CreditCard").Update(&customer).Error

				if dbException == nil {
					context.JSON(http.StatusOK, customer)
				} else {
					respondError(context, http.StatusInternalServerError, dbException.Error())
				}
			} else {
				respondError(context, http.StatusBadRequest, jsonException.Error())
			}
		} else {
			respondError(context, http.StatusUnauthorized, "permissionDenied")
		}
	} else {
		respondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func Delete(context *gin.Context, database *gorm.DB, session string) {

	username, usernameFound := context.Params.Get("id")

	if usernameFound {

		if session == "admin" || session == username {

			dbException := database.Delete(&model.Customer{Username: username}).Error

			if dbException == nil {
				context.JSON(http.StatusNoContent, nil)
			} else {
				respondError(context, http.StatusInternalServerError, dbException.Error())
			}
		} else {
			respondError(context, http.StatusUnauthorized, "permissionDenied")
		}
	} else {
		respondError(context, http.StatusBadRequest, "missingParameter")
	}
}

func GetCustomer(database *gorm.DB, username string) *model.Customer {

	customer := model.Customer{}
	dbException := database.First(&customer, model.Customer{Username: username}).Error

	if dbException == nil {
		return &customer
	} else {
		return nil
	}
}

func GetCustomerOr404(context *gin.Context, database *gorm.DB, username string) *model.Customer {

	customer := model.Customer{}
	dbException := database.Preload("CreditCard").First(&customer, model.Customer{Username: username}).Error

	if dbException == nil {
		return &customer
	} else {
		respondError(context, http.StatusNotFound, dbException.Error())
	}

	return nil
}

func respondError(context *gin.Context, statusCode int, message string) {
	context.JSON(statusCode, gin.H{"error": message})
}
