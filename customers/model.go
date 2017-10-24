package customers

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

type CreditCard struct {
	ID       int       `json:"id" gorm:"primary_key"`
	Type     string    `json:"type" gorm:"not null"`
	Number   string    `json:"number" gorm:"not null"`
	Validity time.Time `json:"validity" gorm:"not null"`
}

type Customer struct {
	common.Model
	Name         string      `json:"name" gorm:"not null"`
	Username     string      `json:"username" gorm:"not null"`
	PublicKey    string      `json:"key" gorm:"not null"`
	Password     string      `json:"password" gorm:"not null"`
	Email        string      `json:"email" gorm:"not null;unique"`
	Address      string      `json:"address" gorm:"not null"`
	TaxNumber    string      `json:"nif" sql:"gorm null;size:9"`
	Country      string      `json:"country" gorm:"not null;size:2"`
	CreditCard   *CreditCard `json:"credit_card"`
	CreditCardID uint        `json:"credit_card_id"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&CreditCard{})
	database.AutoMigrate(&Customer{})
	database.Model(&CreditCard{}).Related(Customer{}, "CreditCardID")
	database.Model(&Customer{}).AddForeignKey("credit_card_id", "credit_cards(id)", "CASCADE", "CASCADE")
}
