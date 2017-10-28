package products

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
)

func GetProductsByBarcode(database *sqlx.DB, barcodeList []string) ([]Product, error) {

	if query, args, errors := preloadBarcode.From(Products).Where(
		squirrel.Eq{Barcode: barcodeList},
	).ToSql(); errors != nil {
		return []Product{}, errors
	} else {
		var products []Product
		return products, database.Select(&products, query, args...)
	}
}

func (product *Product) GenerateJson() ProductJSON {

	return ProductJSON{
		Name:        product.Name,
		Brand:       product.Brand,
		Price:       product.Price,
		Barcode:     product.Barcode,
		ImageUri:    product.ImageUri,
		Description: product.Description,
	}
}