package customers

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/common"
)

func List(database *gorm.DB, username string) (int, interface{}) {

	if username != common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customers := []Customer{}
	database.Preload("CreditCard").Find(&customers)
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
		Password:  hashedPassword,
		Name:      customerPOST.Name,
		Username:  customerPOST.Username,
		Country:   customerPOST.Country,
		Address1:  customerPOST.Address1,
		Address2:  customerPOST.Address2,
		TaxNumber: customerPOST.TaxNumber,
		PublicKey: customerPOST.PublicKey,
		CreditCard: &CreditCard{
			Number:   creditCard.Number,
			Type:     creditCard.Type,
			Validity: creditCard.Validity,
		},
	}

	dbException := database.Save(&customer).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	return http.StatusCreated, generateJson(customer)
}

func Find(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != customerId {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customer := Customer{}
	dbException := database.Preload("CreditCard").First(&customer, "username = ?", customerId).Error

	if dbException != nil {
		return http.StatusNotFound, common.JSON(dbException)
	}

	return http.StatusOK, generateJson(customer)
}

func Update(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != customerId {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customer := Customer{}
	dbException := database.Preload("CreditCard").First(&customer, username).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	jsonException := context.Bind(&customer)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	}

	dbException = database.Preload("CreditCard").Update(&customer).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	return http.StatusOK, generateJson(customer)
}

func Delete(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != customerId {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	dbException := database.Delete(&Customer{Username: customerId}).Error

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	return http.StatusNoContent, nil
}

func Authenticate(database *gorm.DB, username string, password string) (string, bool) {

	customer := Customer{Username: username}

	if dbException := database.First(&customer, &customer).Error; dbException != nil {
		return customer.Username, false
	}

	return customer.Username, auth.VerifyPassword(customer.Password, password) == nil
}