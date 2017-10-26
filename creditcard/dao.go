package creditcard

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
)

const (
	Type        = "type"
	Number      = "number"
	Validity    = "validity"
	CreditCards = "credit_cards"
)

func GetById(database *sqlx.DB, creditCardId int) (*CreditCard, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select("*").From(CreditCards).Where(
		squirrel.Eq{common.Id: creditCardId},
	).Limit(1).ToSql()

	if sqlException != nil {
		return nil, sqlException
	} else {
		var creditCard CreditCard
		return &creditCard, database.Get(&creditCard, sqlQuery, sqlArgs...)
	}
}

func Insert(database *sqlx.DB, creditCardJSON *CreditCardJSON) (*CreditCard, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Insert(CreditCards).Columns(
		Type, Number, Validity,
	).Values(
		creditCardJSON.Type, creditCardJSON.Number, creditCardJSON.Validity,
	).Suffix("RETURNING *").ToSql()

	if sqlException != nil {
		return nil, sqlException
	} else {
		var creditCard CreditCard
		return &creditCard, database.Get(&creditCard, sqlQuery, sqlArgs...)
	}
}

func Update(database *sqlx.DB, creditCardId int, creditCardJSON *CreditCardJSON) (*CreditCard, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Update(CreditCards).SetMap(map[string]interface{}{
		Type:     creditCardJSON.Type,
		Number:   creditCardJSON.Number,
		Validity: creditCardJSON.Validity,
	}).Where(
		squirrel.Eq{common.Id: creditCardId},
	).Suffix("RETURNING *").ToSql()

	if sqlException != nil {
		return nil, sqlException
	} else {
		var creditCard CreditCard
		return &creditCard, database.Get(&creditCard, sqlQuery, sqlArgs...)
	}
}