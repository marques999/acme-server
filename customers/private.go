package customers

import (
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
)

func generateCustomer(query *sqlx.Rows, invokeNext bool) *Customer {

	if invokeNext && query.Next() == false {
		return nil
	}

	customer := Customer{}
	defer query.Close()

	query.Scan(&customer.ID, &customer.Name, &customer.Country,
		&customer.Username, &customer.Password, &customer.Address1,
		&customer.Address2, &customer.PublicKey, &customer.TaxNumber,
		&customer.CreatedAt, &customer.UpdatedAt, &customer.CreditCardID,
		&customer.CreditCard.ID, &customer.CreditCard.Type,
		&customer.CreditCard.Number, &customer.CreditCard.Validity)

	return &customer
}

func (customer *Customer) generateDetails(creditCard *CreditCard) map[string]interface{} {

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
		CreditCardData: CreditCardJSON{
			Type:     creditCard.Type,
			Number:   creditCard.Number,
			Validity: creditCard.Validity,
		},
		common.CreatedAt: customer.CreatedAt,
		common.UpdatedAt: customer.UpdatedAt,
	}
}