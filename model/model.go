package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

func Migrate(database *gorm.DB) *gorm.DB {

	database.AutoMigrate(&CreditCard{})
	database.AutoMigrate(&Product{})
	database.AutoMigrate(&Customer{})
	database.Model(&CreditCard{}).Related(Customer{}, "CreditCardID")
	database.Model(&Customer{}).AddForeignKey("credit_card_id", "credit_cards(id)", "CASCADE", "CASCADE")

	return database
}

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"modified"`
	DeletedAt *time.Time `json:"deleted" gorm:"index"`
}

type Product struct {
	Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Barcode     string  `json:"barcode" gorm:"not null;index;unique_index"`
}

type Customer struct {
	Model
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

type CreditCard struct {
	Model
	Type     string    `json:"type" gorm:"not null"`
	Number   string    `json:"number" gorm:"not null"`
	Validity time.Time `json:"validity" gorm:"not null"`
}
