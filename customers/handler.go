package customers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/auth"
)

func List(database *gorm.DB) (int, interface{}) {
	customers := []Customer{}
	database.Preload("CreditCard").Find(&customers)
	return http.StatusOK, customers
}

func Insert(context *gin.Context, database *gorm.DB) (int, interface{}) {

	customer := Customer{}

	if ex := context.Bind(&customer); ex == nil {
		return http.StatusBadRequest, ex.Error()
	} else if customer.Name == common.AdminAccount {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	password, ex := auth.GeneratePassword(customer.Password)
	customer.Password = password

	if ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := database.Save(&customer).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusCreated, customer
	}
}

func Find(context *gin.Context, database *gorm.DB, session string) (int, interface{}) {

	customer := Customer{}

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if session != common.AdminAccount && session != id {
		return http.StatusUnauthorized, common.PermissionDenied()
	} else if dbException := database.Preload("CreditCard").First(&customer, "username = ?", id).Error; dbException != nil {
		return http.StatusNotFound, common.JSON(dbException)
	} else {
		return http.StatusOK, customer
	}
}

func Update(context *gin.Context, database *gorm.DB, session string) (int, interface{}) {

	customer := Customer{}

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if session != common.AdminAccount && session != id {
		return http.StatusUnauthorized, common.PermissionDenied()
	} else if ex := database.Preload("CreditCard").First(&customer, "username = ?", session).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := context.Bind(&customer); ex != nil {
		return http.StatusBadRequest, common.JSON(ex)
	} else if ex := database.Preload("CreditCard").Update(&customer).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusOK, customer
	}
}

func Delete(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if username != common.AdminAccount && username != id {
		return http.StatusUnauthorized, common.PermissionDenied()
	} else if ex := database.Delete(&Customer{Username: id}).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusNoContent, nil
	}
}
