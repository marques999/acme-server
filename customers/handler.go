package customers

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/common"
)

func List(database *gorm.DB, username string) (int, interface{}) {

	customers := []Customer{}

	if username != common.AdminAccount {
		return http.StatusUnauthorized, nil
	} else {
		database.Preload("CreditCard").Find(&customers)
	}

	jsonCustomers := make([]CustomerJSON, len(customers))

	for i, v := range customers {
		jsonCustomers[i] = generateJson(v)
	}

	return http.StatusOK, jsonCustomers
}

func Insert(context *gin.Context, database *gorm.DB) (int, interface{}) {

	customerPOST := CustomerPOST{}
	jsonException := context.Bind(&customerPOST)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	} else if customerPOST.Username == common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	creditCard := customerPOST.CreditCard
	hashedPassword, hashException := auth.GeneratePassword(customerPOST.Password)

	if hashException != nil {
		return http.StatusInternalServerError, common.JSON(hashException)
	}

	customer := Customer{
		Name:      customerPOST.Name,
		Username:  customerPOST.Username,
		Password:  hashedPassword,
		Country:   customerPOST.Country,
		Address:   customerPOST.Address,
		TaxNumber: customerPOST.TaxNumber,
		PublicKey: customerPOST.PublicKey,
		CreditCard: &CreditCard{
			Number:   creditCard.Number,
			Type:     creditCard.Type,
			Validity: creditCard.Validity,
		},
	}

	dbException := database.Save(&customer).Error

	if dbException == nil {
		return http.StatusCreated, generateJson(customer)
	} else {
		return http.StatusInternalServerError, common.JSON(dbException)
	}
}

func Find(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != customerId {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customer := Customer{Username: customerId}
	dbException := database.Preload("CreditCard").First(&customer, &customer).Error

	if dbException == nil {
		return http.StatusOK, generateJson(customer)
	} else {
		return http.StatusNotFound, common.JSON(dbException)
	}
}

func Update(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != customerId {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customer := Customer{}
	dbsException := database.Preload("CreditCard").First(&customer, username).Error

	if dbsException != nil {
		return http.StatusInternalServerError, common.JSON(dbsException)
	}

	jsonException := context.Bind(&customer)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	}

	dbuException := database.Preload("CreditCard").Update(&customer).Error

	if dbuException == nil {
		return http.StatusOK, generateJson(customer)
	} else {
		return http.StatusInternalServerError, common.JSON(dbuException)
	}
}

func Delete(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != customerId {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	dbException := database.Delete(&Customer{Username: customerId}).Error

	if dbException == nil {
		return http.StatusNoContent, nil
	} else {
		return http.StatusInternalServerError, common.JSON(dbException)
	}
}

func Authenticate(database *gorm.DB, username string, password string) (string, bool) {

	customer := Customer{Username: username}

	if dbException := database.First(&customer, &customer).Error; dbException != nil {
		return customer.Username, false
	} else {
		return customer.Username, auth.VerifyPassword(customer.Password, password) == nil
	}
}