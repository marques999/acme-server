package customers

import (
	"strings"
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/creditcard"
)

func getCustomer(database *sqlx.DB, customerId string) (int, interface{}) {

	customer, sqlException := QueryByUsername(database, customerId)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	}

	creditCard, sqlException := creditcard.GetById(database, customer.CreditCardID)

	if sqlException != nil {
		return http.StatusInternalServerError, common.JSON(sqlException)
	}

	return http.StatusOK, customer.generateJson(creditCard)
}

func (customer *Customer) generateJson(creditCard *creditcard.CreditCard) CustomerJSON {

	if len(creditCard.Number) > 4 {

		creditCard.Number = strings.Repeat(
			"*", len(creditCard.Number)-4,
		) + creditCard.Number[len(creditCard.Number)-4:]
	}

	return CustomerJSON{
		Name:       customer.Name,
		Username:   customer.Username,
		Address1:   customer.Address1,
		Address2:   customer.Address2,
		Country:    customer.Country,
		TaxNumber:  customer.TaxNumber,
		CreditCard: creditCard.GenerateJson(),
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
	}
}