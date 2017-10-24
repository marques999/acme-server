package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
)

type Order struct {
	common.Model
	Token    string
	Customer int                `sql:"not null"`
	Total    float64            `sql:"DEFAULT:0;not null"`
	Status   int                `sql:"DEFAULT:0;not null"`
	Products []products.Product `sql:"save_associations:false;many2many:order_products"`
}

type OrderPOST struct {
	Payload   string `binding:"required" json:"payload"`
	Signature string `binding:"required" json:"signature"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&Order{})
	database.Model(&Order{}).Related(&products.Product{}, "Products")
	database.Model(&Order{}).AddForeignKey("customer", "customers(id)", "CASCADE", "CASCADE")
}