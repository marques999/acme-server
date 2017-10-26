package customers

import (
	"time"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
)

func Query(database *sqlx.DB) ([]Customer, error) {

	customers := []Customer{}
	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select("*").From(Customers).ToSql()

	if sqlException != nil {
		return customers, sqlException
	} else {
		return customers, database.Select(&customers, sqlQuery, sqlArgs...)
	}
}

func QueryByUsername(database *sqlx.DB, username string) (*Customer, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select(
		Username, Name, Country, TaxNumber, Address1, Address2, Password, CreditCard,
	).From(Customers).Where(
		squirrel.Eq{Username: username},
	).Limit(1).ToSql()

	if sqlException != nil {
		return nil, sqlException
	} else {
		var customer Customer
		return &customer, database.Get(&customer, sqlQuery, sqlArgs...)
	}
}

func insertCustomer(database *sqlx.DB, customerPOST CustomerPOST, creditCardId int) (*Customer, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Insert(Customers).Columns(
		Name, Country, Username, Password, Address1, Address2, PublicKey, TaxNumber, CreditCard,
	).Values(
		customerPOST.Name,
		customerPOST.Country,
		customerPOST.Username,
		customerPOST.Password,
		customerPOST.Address1,
		customerPOST.Address2,
		customerPOST.PublicKey,
		customerPOST.TaxNumber,
		creditCardId,
	).Suffix("RETURNING *").ToSql()

	if sqlException != nil {
		return nil, sqlException
	} else {
		var customer Customer
		return &customer, database.Get(&customer, sqlQuery, sqlArgs...)
	}
}

func Update(database *sqlx.DB, customerId string, customerPOST *CustomerPOST) (*Customer, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Update(Customers).SetMap(map[string]interface{}{
		Password:         customerPOST.Password,
		Address1:         customerPOST.Address1,
		Address2:         customerPOST.Address2,
		TaxNumber:        customerPOST.TaxNumber,
		Country:          customerPOST.Country,
		Name:             customerPOST.Name,
		common.UpdatedAt: time.Now(),
	}).Where(
		squirrel.Eq{Username: customerId},
	).Suffix("RETURNING *").ToSql()

	if sqlException != nil {
		return nil, sqlException
	} else {
		var customer Customer
		return &customer, database.Get(&customer, sqlQuery, sqlArgs...)
	}
}

func deleteCustomer(db *sqlx.DB, username string) (sql.Result, error) {

	return common.StatementBuilder().Delete(Customers).Where(
		squirrel.Eq{Username: username},
	).RunWith(db).Exec()
}
