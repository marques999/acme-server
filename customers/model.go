package customers

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

type CreditCard struct {
	ID       int       `sql:"primary_key"`
	Type     string    `sql:"not null"`
	Number   string    `sql:"not null"`
	Validity time.Time `sql:"not null"`
}

type CreditCardJSON struct {
	Type     string    `binding:"required" json:"type"`
	Number   string    `binding:"required" json:"number"`
	Validity time.Time `binding:"required" json:"validity"`
}

type Customer struct {
	common.Model
	Name         string `sql:"not null"`
	PublicKey    string `sql:"not null"`
	Password     string `sql:"not null"`
	Username     string `sql:"not null;unique;unique_text"`
	Address      string `sql:"not null"`
	TaxNumber    string `sql:"gorm null;size:9"`
	Country      string `sql:"not null;size:2"`
	CreditCardID int
	CreditCard   *CreditCard
}

type CustomerJSON struct {
	Name       string         `json:"name"`
	Username   string         `json:"email"`
	Address    string         `json:"address"`
	TaxNumber  string         `json:"nif"`
	Country    string         `json:"country"`
	CreatedAt  time.Time      `json:"created"`
	UpdatedAt  time.Time      `json:"modified"`
	CreditCard CreditCardJSON `json:"credit_card"`
}

type CustomerPOST struct {
	Name       string         `binding:"required" json:"name"`
	PublicKey  string         `binding:"required" json:"key"`
	Password   string         `binding:"required" json:"password"`
	Username   string         `binding:"required" json:"email"`
	Address    string         `binding:"required" json:"address"`
	TaxNumber  string         `binding:"required" json:"nif"`
	Country    string         `binding:"required" json:"country"`
	CreditCard CreditCardJSON `binding:"required" json:"credit_card"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&CreditCard{})
	database.AutoMigrate(&Customer{})
	database.Model(&CreditCard{}).Related(Customer{}, "CreditCardID")
	database.Model(&Customer{}).AddForeignKey("credit_card_id", "credit_cards(id)", "CASCADE", "CASCADE")
}