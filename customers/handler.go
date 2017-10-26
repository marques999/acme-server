package customers

import (
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/creditcard"
	"github.com/Masterminds/squirrel"
)

func LIST(database *sqlx.DB, username string) (int, interface{}) {

	if username != common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customers, sqlException := Query(database)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	}

	jsonCustomers := make([]CustomerJSON, len(customers))

	for index, customer := range customers {

		creditCard, sqlException := creditcard.GetById(database, customer.CreditCardID)

		if sqlException == nil {
			jsonCustomers[index] = customer.generateJson(creditCard)
		} else {
			return http.StatusInternalServerError, common.JSON(sqlException)
		}
	}

	return http.StatusOK, jsonCustomers
}

func GET(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if common.HasPermissions(username, customerId) == false {
		return http.StatusUnauthorized, common.PermissionDenied()
	} else {
		return getCustomer(database, customerId)
	}
}

func POST(context *gin.Context, database *sqlx.DB) (int, interface{}) {

	customerPOST := CustomerPOST{}
	jsonException := context.Bind(&customerPOST)

	if jsonException != nil {
		return http.StatusBadRequest, common.JSON(jsonException)
	} else if customerPOST.Username == common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	creditCard, dbException := creditcard.Insert(database, &customerPOST.CreditCard)

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	}

	hashedPassword, hashException := auth.GeneratePassword(customerPOST.Password)

	if hashException != nil {
		return http.StatusInternalServerError, common.JSON(hashException)
	}

	customerPOST.Password = hashedPassword
	customer, sqlException := insertCustomer(database, customerPOST, creditCard.ID)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	} else {
		return http.StatusOK, customer.generateJson(creditCard)
	}
}

func PUT(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if common.HasPermissions(username, customerId) == false {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	customerPost := CustomerPOST{}
	bindException := context.Bind(&customerPost)

	if bindException != nil {
		return http.StatusBadRequest, common.JSON(bindException)
	}

	customer, sqlException := Update(database, customerId, &customerPost)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	}

	creditCard, sqlException := creditcard.GetById(database, customer.CreditCardID)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	}

	return http.StatusOK, customer.generateJson(creditCard)
}

func DELETE(context *gin.Context, database *sqlx.DB, username string) (int, interface{}) {

	customerId, paramExists := context.Params.Get("id")

	if paramExists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if common.HasPermissions(username, customerId) == false {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	_, dbException := deleteCustomer(database, username)

	if dbException != nil {
		return http.StatusInternalServerError, common.JSON(dbException)
	} else {
		return http.StatusNoContent, nil
	}
}

func Authenticate(database *sqlx.DB, username string, password string) (string, bool) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select(Password).From(Customers).Where(
		squirrel.Eq{Username: username},
	).Limit(1).ToSql()

	if sqlException != nil {
		return username, false
	}

	customer := Customer{}
	sqlException = database.Get(&customer, sqlQuery, sqlArgs...)

	if sqlException != nil {
		return username, false
	} else {
		return username, auth.VerifyPassword(customer.Password, password) == nil
	}
}