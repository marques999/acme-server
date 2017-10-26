package auth

import (
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/common"
)

func Clean(database *sqlx.DB) (int, interface{}) {

	if _, ex := common.StatementBuilder().Delete(
		"products",
	).Suffix("CASCADE").RunWith(database).Exec(); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if _, ex := common.StatementBuilder().Delete(
		"customers",
	).Suffix("CASCADE").RunWith(database).Exec(); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, nil
	}
}