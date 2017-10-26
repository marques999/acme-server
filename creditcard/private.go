package creditcard

func (creditCard *CreditCard) GenerateJson() *CreditCardJSON {

	return &CreditCardJSON{
		Type:     creditCard.Type,
		Number:   creditCard.Number,
		Validity: creditCard.Validity,
	}
}
