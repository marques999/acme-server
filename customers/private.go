package customers

import "github.com/jmoiron/sqlx"

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