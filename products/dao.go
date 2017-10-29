package products

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
)

var preloadList = common.SqlBuilder().Select("*").From(Products)
var preloadGet = common.SqlBuilder().Select("*").From(Products).Limit(1)

func getProducts(database *sqlx.DB) ([]Product, error) {

	if query, args, errors := preloadList.ToSql(); errors != nil {
		return []Product{}, errors
	} else {
		var products []Product
		return products, database.Select(&products, query, args...)
	}
}

func getProduct(database *sqlx.DB, barcode string) (*Product, error) {

	if query, args, errors := preloadGet.Where(
		squirrel.Eq{Barcode: barcode},
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var product Product
		return &product, database.Get(&product, query, args...)
	}
}

var preloadInsert = common.SqlBuilder().Insert(Products).Columns(
	Name, Brand, Price, Barcode, ImageUri, Description,
).Suffix(common.ReturningRow)

func insertProduct(database *sqlx.DB, productJson ProductJSON) (*Product, error) {

	if query, args, errors := preloadInsert.Values(
		productJson.Name,
		productJson.Brand,
		productJson.Price,
		productJson.Barcode,
		productJson.ImageUri,
		productJson.Description,
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var product Product
		return &product, database.Get(&product, query, args...)
	}
}

var preloadBarcode = common.SqlBuilder().Select("*")
var preloadDelete = common.SqlBuilder().Delete(Products)
var preloadUpdate = common.SqlBuilder().Update(Products).Suffix(common.ReturningRow)

func updateProduct(database *sqlx.DB, barcode string, productJson ProductJSON) (*Product, error) {

	if query, args, errors := preloadUpdate.SetMap(map[string]interface{}{
		Name:        productJson.Name,
		Brand:       productJson.Brand,
		Price:       productJson.Price,
		ImageUri:    productJson.ImageUri,
		Description: productJson.Description,
	}).Where(squirrel.Eq{Barcode: barcode}).ToSql(); errors != nil {
		return nil, errors
	} else {
		var product Product
		return &product, database.Get(&product, query, args...)
	}
}

func deleteProduct(database *sqlx.DB, barcode string) (sql.Result, error) {
	return preloadDelete.Where(squirrel.Eq{Barcode: barcode}).RunWith(database.DB).Exec()
}