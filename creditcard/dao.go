package creditcard

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
)

var preloadGet = common.SqlBuilder().Select("*").From(CreditCards).Limit(1)
var preloadUpdate = common.SqlBuilder().Update(CreditCards).Suffix(common.ReturningRow)

func GetById(database *sqlx.DB, creditCardId int) (*CreditCard, error) {

	if query, args, errors := preloadGet.Where(
		squirrel.Eq{common.Id: creditCardId},
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var creditCard CreditCard
		return &creditCard, database.Get(&creditCard, query, args...)
	}
}

var preloadInsert = common.SqlBuilder().Insert(CreditCards).Columns(
	Type, Number, Validity,
).Suffix(common.ReturningRow)

func Insert(database *sqlx.DB, creditCardJSON *CreditCardJSON) (*CreditCard, error) {

	if query, args, errors := preloadInsert.Values(
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

func Update(database *sqlx.DB, creditCardId int, creditCardJSON *CreditCardJSON) (*CreditCard, error) {

	if query, args, errors := preloadUpdate.SetMap(map[string]interface{}{
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