package products

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
)

func getProducts(database *sqlx.DB) ([]Product, error) {

	products := []Product{}
	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select("*").From(Products).ToSql()

	if sqlException != nil {
		return products, sqlException
	} else {
		return products, database.Select(&products, sqlQuery, sqlArgs...)
	}
}

func GetProductsByBarcode(database *sqlx.DB, barcodes []string) ([]Product, error) {

	products := []Product{}

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select("*").From(Products).Where(
		squirrel.Eq{Barcode: barcodes},
	).ToSql()

	if sqlException != nil {
		return products, sqlException
	} else {
		return products, database.Select(&products, sqlQuery, sqlArgs...)
	}
}

func getProduct(database *sqlx.DB, barcode string) (*Product, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select("*").From(Products).Where(
		squirrel.Eq{Barcode: barcode},
	).Limit(1).ToSql()

	product := Product{}

	if sqlException != nil {
		return nil, sqlException
	} else {
		return &product, database.Get(&product, sqlQuery, sqlArgs...)
	}
}

func insertProduct(database *sqlx.DB, productJson ProductJSON) (*Product, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Insert(Products).Columns(
		Name, Brand, Price, Barcode, ImageUri, Description,
	).Values(
		productJson.Name,
		productJson.Brand,
		productJson.Price,
		productJson.Barcode,
		productJson.ImageUri,
		productJson.Description,
	).Suffix("RETURNING *").ToSql()

	product := Product{}

	if sqlException == nil {
		return &product, database.Get(&product, sqlQuery, sqlArgs...)
	} else {
		return nil, sqlException
	}
}

func updateProduct(db *sqlx.DB, productJson ProductJSON) (*Product, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Update(Products).SetMap(map[string]interface{}{
		Name:        productJson.Name,
		Brand:       productJson.Brand,
		Price:       productJson.Price,
		ImageUri:    productJson.ImageUri,
		Description: productJson.Description,
	}).Suffix("RETURNING *").ToSql()

	product := Product{}

	if sqlException == nil {
		return nil, sqlException
	} else {
		return &product, db.Select(&product, sqlQuery, sqlArgs...)
	}
}

func deleteProduct(db *sqlx.DB, barcode string) (sql.Result, error) {
	return common.StatementBuilder().Delete(Products).Where(squirrel.Eq{Barcode: barcode}).RunWith(db).Exec()
}