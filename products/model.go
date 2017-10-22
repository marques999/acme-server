package products

import (
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/common"
)

type Product struct {
	common.Model
	Name        string  `json:"name" gorm:"not null"`
	Brand       string  `json:"brand" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Barcode     string  `json:"barcode" gorm:"not null;index;unique_index"`
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&Product{})
}
