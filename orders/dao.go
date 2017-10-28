package orders

import (
	"time"
	"math/rand"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
	"github.com/marques999/acme-server/customers"
)

var preloadGet = common.SqlBuilder().Select(
	"orders.id",
	"COUNT(*)",
	"SUM(price * quantity) AS total",
	"orders.created_at",
	"orders.updated_at",
	"orders.customer",
	"orders.status",
	"orders.token",
).From(Orders).Join(
	"order_products ON orders.id = order_products.order_id",
).Join(
	"products ON products.barcode = order_products.product_id",
).GroupBy("orders.id")

func listOrders(database *sqlx.DB, username string) ([]Order, error) {

	builder := preloadGet

	if username != common.AdminAccount {
		builder = builder.Where(squirrel.Eq{Customer: username})
	}

	if query, args, errors := builder.ToSql(); errors != nil {
		return []Order{}, nil
	} else {
		var orders []Order
		return orders, database.Select(&orders, query, args...)
	}
}

func getOrder(database *sqlx.DB, condition squirrel.Eq) (*Order, error) {

	if query, args, errors := preloadGet.Where(condition).ToSql(); errors != nil {
		return nil, errors
	} else {
		var order Order
		return &order, database.Get(&order, query, args...)
	}
}

var preloadDelete = common.SqlBuilder().Delete(Orders)
var preloadUpdate = common.SqlBuilder().Update(Orders).Suffix(common.ReturningRow)
var preloadManyDelete = common.SqlBuilder().Delete(OrderProducts)

func updateOrder(database *sqlx.DB, condition squirrel.Eq, status int) (*Order, error) {

	if query, args, errors := preloadUpdate.SetMap(map[string]interface{}{
		Status: status,
	}).Where(condition).ToSql(); errors != nil {
		return nil, errors
	} else {
		var order Order
		return &order, database.Get(&order, query, args...)
	}
}

var preloadInsert = common.SqlBuilder().Insert(Orders).Columns(
	Customer, Status,
).Suffix(common.ReturningRow)

func insertOrder(
	database *sqlx.DB,
	customer *customers.Customer,
	customerCartPOST ...CustomerCartPOST,
) (*map[string]interface{}, error) {

	// inserir encomenda sem token nem total

	order := Order{Status: ValidationFailed}
	query, args, errors := preloadInsert.Values(customer, ValidationFailed).ToSql()

	if errors != nil {
		return nil, errors
	}

	// obtém encomenda recentemente inserida

	if errors = database.Get(&order, query, args...); errors != nil {
		return nil, errors
	}

	// inserir as entidades da associação muitos-para-muitos

	customerCart, errors := insertOrderProducts(database, order.ID, customerCartPOST)

	if errors != nil {
		return nil, errors
	}

	if customer.CreditCard.Validity.After(time.Now()) && rand.Float64() <= common.SuccessProbability {
		order.Status = ValidationComplete
	}

	// atualizar encomenda com total a pagar e token gerado aleatoriamente

	if token, errors := order.generateToken(); errors != nil {
		return nil, errors
	} else if query, args, errors = preloadUpdate.SetMap(map[string]interface{}{
		Token:  token,
		Status: order.Status,
	}).Where(squirrel.Eq{common.Id: order.ID}).ToSql(); errors != nil {
		return nil, errors
	} else if errors = database.Get(&order, query, args...); errors != nil {
		return nil, errors
	} else {
		return order.generateJson(customerCart), nil
	}
}

var preloadManyGet = common.SqlBuilder().Select(
	Quantity, "products.*",
).From(OrderProducts).Join(
	"products ON products.barcode = order_products.product_id",
)

func getCustomerCart(database *sqlx.DB, orderId int) ([]CustomerCartJSON, error) {

	if query, args, errors := preloadManyGet.Where(
		squirrel.Eq{OrderID: orderId},
	).ToSql(); errors != nil {
		return []CustomerCartJSON{}, errors
	} else if query, errors := database.Queryx(query, args); errors != nil {
		return []CustomerCartJSON{}, errors
	} else {
		return generateCustomerCart(query), nil
	}
}

var preloadManyInsert = common.SqlBuilder().Insert(OrderProducts).Columns(
	OrderID, ProductID, Quantity,
)

func insertOrderProducts(database *sqlx.DB, orderId int, customerCartPOST []CustomerCartPOST) ([]CustomerCartJSON, error) {

	barcodes := make([]string, len(customerCartPOST))

	for index, entry := range customerCartPOST {
		barcodes[index] = entry.Product
	}

	customerCart := make([]CustomerCartJSON, len(customerCartPOST))
	products, errors := products.GetProductsByBarcode(database, barcodes)

	if errors != nil {
		return customerCart, errors
	}

	builder := preloadManyInsert

	for index, product := range products {

		quantity := customerCartPOST[index].Quantity
		builder = builder.Values(orderId, product.ID, quantity)

		customerCart[index] = CustomerCartJSON{
			Product:  product.GenerateJson(),
			Quantity: quantity,
		}
	}

	if _, errors = builder.RunWith(database.DB).Exec(); errors != nil {
		return []CustomerCartJSON{}, errors
	} else {
		return customerCart, nil
	}
}

func deleteOrder(database *sqlx.DB, condition squirrel.Eq) (sql.Result, error) {
	return preloadDelete.Where(condition).RunWith(database.DB).Exec()
}

func deleteOrderProducts(database *sqlx.DB, orderId int) (sql.Result, error) {
	return preloadManyDelete.Where(squirrel.Eq{OrderID: orderId}).RunWith(database.DB).Exec()
}