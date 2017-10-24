package products

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

type Product struct {
	common.Model
	Name        string  `sql:"not null"`
	Brand       string  `sql:"not null"`
	Price       float64 `sql:"not null"`
	ImageUri    string  `sql:"not null"`
	Description string  `sql:"not null"`
	Barcode     string  `sql:"not null;index;unique_index"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&Product{})
}