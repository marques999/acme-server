package common

import (
	"net/http"
	"github.com/jmoiron/sqlx"
)

var cleanProducts = SqlBuilder().Delete("products").Suffix("CASCADE")
var cleanCustomers = SqlBuilder().Delete("customers").Suffix("CASCADE")

func Clean(database *sqlx.DB) (int, interface{}) {

	if _, errors := cleanProducts.RunWith(database.DB).Exec(); errors != nil {
		return http.StatusInternalServerError, JSON(errors)
	} else if _, errors := cleanCustomers.RunWith(database.DB).Exec(); errors != nil {
		return http.StatusInternalServerError, JSON(errors)
	} else {
		return http.StatusOK, nil
	}
}