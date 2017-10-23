package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/products"
)

type Order struct {
	common.Model
	Valid    bool                `json:"valid"`
	Token    string              `json:"token"`
	Customer *customers.Customer `json:"customer" sql:"save_associations:false;not null"`
	Products []products.Product  `json:"products" sql:"save_associations:false;many2many:order_products"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&Order{})
	database.Model(&Order{}).Related(&products.Product{}, "Products")
}
