package customers

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/auth"
)

func Authenticate(database *gorm.DB, username string, password string) (string, bool) {

	customer := Customer{}

	if dbException := database.Preload("CreditCard").First(&customer, "username = ?", username).Error; dbException != nil {
		return username, false
	} else {
		return username, auth.VerifyPassword(customer.Password, username, password) == nil
	}
}
