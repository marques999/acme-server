package tokens

import (
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
	"github.com/jinzhu/gorm"
)

type Token struct {
	common.Model
	Customer string             `json:"customer" gorm:"not null"`
	Products []products.Product `json:"products" gorm:"save_associations:false;many2many:token_products"`
}

type TokenPOST struct {
	Customer string   `json:"customer"`
	Products []string `json:"products"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&Token{})
	database.Model(&Token{}).Related(&products.Product{}, "Products")
}
