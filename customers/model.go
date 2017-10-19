package customers

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"time"
)

type CreditCard struct {
	common.Model
	Type     string    `json:"type" gorm:"not null"`
	Number   string    `json:"number" gorm:"not null"`
	Validity time.Time `json:"validity" gorm:"not null"`
}

type Customer struct {
	common.Model
	Name         string      `json:"name" gorm:"type:varchar(255)"`
	Password     string      `json:"password" gorm:"not null; type:varchar(32)"`
	Username     string      `json:"username" gorm:"not null; unique; type:varchar(32)"`
	Email        string      `json:"email" gorm:"not null; unique; type:varchar(255)"`
	Address      string      `json:"address" gorm:"not null"`
	TaxNumber    string      `json:"nif" gorm:"not null;size:9"`
	Country      string      `json:"country" gorm:"not null;size:2"`
	CreditCard   *CreditCard `json:"credit_card"`
	CreditCardID uint        `json:"credit_card_id"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&CreditCard{})
	database.AutoMigrate(&Customer{})
	database.Model(&CreditCard{}).Related(Customer{}, "CreditCardID")
	//Model(&Customer{}).AddForeignKey("credit_card_id", "credit_cards(id)", "CASCADE", "CASCADE")
}
