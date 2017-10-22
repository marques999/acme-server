package orders

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/products"
	"github.com/pborman/uuid"
)

func List(database *gorm.DB, customerId string) (int, interface{}) {

	orders := []Order{}

	if customerId == common.AdminAccount {
		database.Preload("Products").Find(&orders)
	} else {
		database.Preload("Products").Find(&orders, "customer = ?", customerId)
	}

	return http.StatusOK, orders
}

func Validate(context *gin.Context, database *gorm.DB, username string) (int, interface{}) {

	order := Order{}

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if ex := database.Preload("Products").First(&order, "id = ?", id).Error; ex != nil {
		return http.StatusNotFound, common.JSON(ex)
	} else if username != common.AdminAccount && username != order.Customer.Username {
		return http.StatusUnauthorized, common.PermissionDenied()
	}

	order.Valid = true
	order.Token = uuid.NewUUID()

	if ex := database.Save(&order).Error; ex != nil {
		return http.StatusNotFound, common.JSON(ex)
	} else {
		return http.StatusOK, order
	}
}

func Insert(context *gin.Context, database *gorm.DB) (int, interface{}) {

	orderPOST := OrderPOST{}
	customer := &customers.Customer{}

	if ex := context.Bind(&orderPOST); ex != nil {
		return http.StatusBadRequest, common.JSON(ex)
	} else if ex := database.Preload("CreditCard").First(&customer, "username = ?", orderPOST.Customer).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	}

	array := []products.Product{}
	order := Order{Customer: customer}

	if ex := database.Where("barcode in (?)", orderPOST.Products).Find(&array).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else if ex := database.Create(&order).Association("Products").Append(array).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusCreated, order
	}
}

func Find(context *gin.Context, database *gorm.DB, customerId string) (int, interface{}) {

	order := Order{}

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if ex := database.Preload("Products").First(&order, getQueryOptions(id, customerId)).Error; ex != nil {
		return http.StatusNotFound, common.JSON(ex)
	} else {
		return http.StatusOK, order
	}
}

func Delete(context *gin.Context, database *gorm.DB, customerId string) (int, interface{}) {

	if id, exists := context.Params.Get("id"); exists == false {
		return http.StatusBadRequest, common.MissingParameter()
	} else if ex := database.Delete(&Order{}, getQueryOptions(id, customerId)).Error; ex != nil {
		return http.StatusInternalServerError, common.JSON(ex)
	} else {
		return http.StatusNoContent, nil
	}
}
