package customers

import (
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/creditcard"
)

func generateCustomer(query *sqlx.Rows, invokeNext bool) *Customer {

	customer := Customer{}
	creditCard := creditcard.CreditCard{}

	if invokeNext && query.Next() == false {
		return nil
	}

	query.Scan(
		&customer.ID, &customer.Name, &customer.Country,
		&customer.Username, &customer.Password, &customer.Address1,
		&customer.Address2, &customer.PublicKey, &customer.TaxNumber,
		&customer.CreatedAt, &customer.UpdatedAt, &customer.CreditCardID,
		&creditCard.ID, &creditCard.Type, &creditCard.Number, &creditCard.Validity,
	)

	defer query.Close()

	if len(creditCard.Number) > 4 {

		creditCard.Number = strings.Repeat(
			"*", len(creditCard.Number)-4,
		) + creditCard.Number[len(creditCard.Number)-4:]
	}

	customer.CreditCard = creditCard

	return &customer
}

func (customer *Customer) generateDetails(creditCard *creditcard.CreditCard) map[string]interface{} {

	if len(creditCard.Number) > 4 {

		creditCard.Number = strings.Repeat(
			"*", len(creditCard.Number)-4,
		) + creditCard.Number[len(creditCard.Number)-4:]
	}

	return gin.H{
		Name:      customer.Name,
		Username:  customer.Username,
		Address1:  customer.Address1,
		Address2:  customer.Address2,
		Country:   customer.Country,
		TaxNumber: customer.TaxNumber,
		CreditCard: creditcard.CreditCardJSON{
			Type:     creditCard.Type,
			Number:   creditCard.Number,
			Validity: creditCard.Validity,
		},
		common.CreatedAt: customer.CreatedAt,
		common.UpdatedAt: customer.UpdatedAt,
	}
}