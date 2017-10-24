package admin

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

func Clean(database *gorm.DB) (int, interface{}) {

	if ex := clearTables(database); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, nil
	}
}

func Populate(database *gorm.DB) (int, interface{}) {

	if ex := clearTables(database); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := populateTables(database); ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, nil
	}
}