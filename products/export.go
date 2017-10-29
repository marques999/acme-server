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
		product.Name, product.Brand,
		product.Price, product.Barcode,
		product.ImageUri, product.Description,
	}
}