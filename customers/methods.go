package customers

import "strings"

func maskNumber(cardNumber string) string {

	if len(cardNumber) <= 4 {
		return cardNumber
	} else {
		return strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
	}
}

func generateJson(customer Customer) CustomerJSON {

	creditCard := customer.CreditCard

	return CustomerJSON{
		Name:      customer.Name,
		Username:  customer.Username,
		Address:   customer.Address,
		Country:   customer.Country,
		TaxNumber: customer.TaxNumber,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
		CreditCard: CreditCardJSON{
			Type:     creditCard.Type,
			Validity: creditCard.Validity,
			Number:   maskNumber(creditCard.Number),
		},
	}
}