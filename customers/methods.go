package customers

import "strings"

func generateJson(customer Customer) CustomerJSON {

	cardNumber := customer.CreditCard.Number

	if len(cardNumber) > 4 {
		cardNumber = strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
	}

	return CustomerJSON{
		Name:      customer.Name,
		Username:  customer.Username,
		Address1:  customer.Address1,
		Address2:  customer.Address2,
		Country:   customer.Country,
		TaxNumber: customer.TaxNumber,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
		CreditCard: CreditCardJSON{
			Number:   cardNumber,
			Type:     customer.CreditCard.Type,
			Validity: customer.CreditCard.Validity,
		},
	}
}