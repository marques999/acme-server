package customers

import (
	"time"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
)

var preloadList = common.SqlBuilder().Select(
	Name, Username, Address1, Address2, TaxNumber,
	common.CreatedAt, common.UpdatedAt,
).From(Customers)

func getCustomers(database *sqlx.DB) ([]CustomerList, error) {

	if query, args, errors := preloadList.ToSql(); errors != nil {
		return []CustomerList{}, errors
	} else {
		var customers []CustomerList
		return customers, database.Select(&customers, query, args...)
	}
}

var preloadGet = common.SqlBuilder().Select("*").From(Customers).Join(
	"credit_cards ON credit_cards.id = customers.credit_card_id",
).Limit(1)

func GetCustomer(database *sqlx.DB, username string) (*Customer, error) {

	if query, args, errors := preloadGet.Where(squirrel.Eq{Username: username}).ToSql(); errors != nil {
		return nil, errors
	} else if result, errors := database.Queryx(query, args...); errors != nil {
		return nil, errors
	} else {
		return generateCustomer(result, true), nil
	}
}

var preloadInsert = common.SqlBuilder().Insert(Customers).Columns(
	Name, Username, Password, Address1,
	Address2, PublicKey, TaxNumber, CreditCardID,
).Suffix(common.ReturningRow)

func insertCustomer(database *sqlx.DB, customerPOST CustomerInsert, creditCardId int) (*Customer, error) {

	if password, errors := common.GeneratePassword(customerPOST.Password); errors != nil {
		return nil, errors
	} else if query, args, errors := preloadInsert.Values(
		customerPOST.Name,
		customerPOST.Username,
		password,
		customerPOST.Address1,
		customerPOST.Address2,
		customerPOST.PublicKey,
		customerPOST.TaxNumber,
		creditCardId,
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var customer Customer
		return &customer, database.Get(&customer, query, args...)
	}
}

var preloadDelete = common.SqlBuilder().Delete(Customers).Suffix(common.ReturningRow)
var preloadUpdate = common.SqlBuilder().Update(Customers).Suffix(common.ReturningRow)
var preloadLogin = common.SqlBuilder().Select(Password).From(Customers).Limit(1)

func updateCustomer(database *sqlx.DB, username string, customerPOST *CustomerUpdate) (*Customer, error) {

	if password, errors := common.GeneratePassword(customerPOST.Password); errors != nil {
		return nil, errors
	} else if query, args, errors := preloadUpdate.SetMap(map[string]interface{}{
		Password:         password,
		Name:             customerPOST.Name,
		Address1:         customerPOST.Address1,
		Address2:         customerPOST.Address2,
		TaxNumber:        customerPOST.TaxNumber,
		common.UpdatedAt: time.Now(),
	}).Where(squirrel.Eq{Username: username}).ToSql(); errors != nil {
		return nil, errors
	} else {
		var customer Customer
		return &customer, database.Get(&customer, query, args...)
	}
}

func validateLogin(database *sqlx.DB, username string) (string, error) {

	if query, args, errors := preloadLogin.Where(
		squirrel.Eq{Username: username},
	).ToSql(); errors != nil {
		return "", errors
	} else {
		var password string
		return password, database.Get(&password, query, args...)
	}
}

func deleteCustomer(database *sqlx.DB, username string) (*Customer, error) {

	if query, args, errors := preloadDelete.Where(
		squirrel.Eq{Username: username},
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var customer Customer
		return &customer, database.Get(&customer, query, args...)
	}
}

var preloadCardInsert = common.SqlBuilder().Insert(CreditCards).Columns(
	Type, Number, Validity,
).Suffix(common.ReturningRow)

func insertCreditCard(database *sqlx.DB, creditCardJSON *CreditCardJSON) (*CreditCard, error) {

	if query, args, errors := preloadCardInsert.Values(
		creditCardJSON.Type,
		creditCardJSON.Number,
		creditCardJSON.Validity,
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var creditCard CreditCard
		return &creditCard, database.Get(&creditCard, query, args...)
	}
}

var preloadCardDelete = common.SqlBuilder().Delete(CreditCards)
var preloadCardUpdate = common.SqlBuilder().Update(CreditCards).Suffix(common.ReturningRow)

func updateCreditCard(database *sqlx.DB, creditCardId int, creditCardJSON *CreditCardJSON) (*CreditCard, error) {

	if query, args, errors := preloadCardUpdate.SetMap(map[string]interface{}{
		Type:     creditCardJSON.Type,
		Number:   creditCardJSON.Number,
		Validity: creditCardJSON.Validity,
	}).Where(squirrel.Eq{common.Id: creditCardId}).ToSql(); errors != nil {
		return nil, errors
	} else {
		var creditCard CreditCard
		return &creditCard, database.Get(&creditCard, query, args...)
	}
}

func deleteCreditCard(database *sqlx.DB, creditCardID int) (sql.Result, error) {
	return preloadCardDelete.Where(squirrel.Eq{common.Id: creditCardID}).RunWith(database.DB).Exec()
}