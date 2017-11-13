package customers

import (
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/marques999/acme-server/common"
)

func (customer *Customer) GenerateDetails(creditCard *CreditCard) map[string]interface{} {

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