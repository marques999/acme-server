package customers

import (
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/marques999/acme-server/common"
)

func Login(database *sqlx.DB, username string, original string) (string, bool) {

	if hashed, errors := validateLogin(database, username); errors != nil {
		return username, false
	} else {
		return username, bcrypt.CompareHashAndPassword([]byte(hashed), []byte(original)) == nil
	}
}

func List(database *sqlx.DB, username string) (int, interface{}) {

	if username != common.AdminAccount {
		return common.PermissionDenied()
	} else if customers, errors := getCustomers(database); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, customers
	}
}

func Find(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	if id, exists := context.Params.Get(common.Id); exists == false {
		return common.MissingParameter()
	} else if common.HasPermissions(username, id) == false {
		return common.PermissionDenied()
	} else if customer, errors := GetCustomer(database, id); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, customer.GenerateDetails(&customer.CreditCard)
	}
}

func Post(context *gin.Context, database *sqlx.DB) (int, interface{}) {

	customerPOST := CustomerInsert{}

	if errors := context.Bind(&customerPOST); errors != nil {
		return http.StatusBadRequest, common.JSON(errors)
	} else if customerPOST.Username == common.AdminAccount {
		return common.PermissionDenied()
	} else if creditCard, errors := insertCreditCard(database, &customerPOST.CreditCard); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else if customer, errors := insertCustomer(database, customerPOST, creditCard.ID); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, customer.GenerateDetails(creditCard)
	}
}

func Put(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	customerPOST := CustomerUpdate{}

	if id, exists := context.Params.Get(common.Id); exists == false {
		return common.MissingParameter()
	} else if common.HasPermissions(username, id) == false {
		return common.PermissionDenied()
	} else if errors := context.Bind(&customerPOST); errors != nil {
		return http.StatusBadRequest, common.JSON(errors)
	} else if customer, errors := updateCustomer(database, id, &customerPOST); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else if creditCard, errors := updateCreditCard(
		database, customer.CreditCardID, &customerPOST.CreditCard,
	); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusOK, customer.GenerateDetails(creditCard)
	}
}

func Delete(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	if id, exists := context.Params.Get(common.Id); exists == false {
		return common.MissingParameter()
	} else if common.HasPermissions(username, id) == false {
		return common.PermissionDenied()
	} else if customer, errors := deleteCustomer(database, id); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else if _, errors := deleteCreditCard(database, customer.CreditCardID); errors != nil {
		return http.StatusInternalServerError, common.JSON(errors)
	} else {
		return http.StatusNoContent, nil
	}
}