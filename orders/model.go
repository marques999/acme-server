package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/products"
	"github.com/pborman/uuid"
)

type Order struct {
	common.Model
	Valid    bool                `json:"valid"`
	Token    uuid.UUID           `json:"token" gorm:"uuid"`
	Customer *customers.Customer `json:"customer" sql:"not null"`
	Products []products.Product  `json:"products" sql:"save_associations:false;many2many:order_products"`
}

type OrderPOST struct {
	Customer string   `json:"customer"`
	Products []string `json:"products"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&Order{})
	database.Model(&Order{}).Related(&products.Product{}, "Products")
}
