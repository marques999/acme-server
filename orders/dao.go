package orders

import (
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

func getOrder(database *sqlx.DB, token string, customer string) (*Order, error) {

	if query, args, errors := preloadGet.Where(getQueryOptions(
		token, customer,
	)).ToSql(); errors != nil {
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
	customerCartPOST []CustomerCartPOST,
) (*map[string]interface{}, error) {

	order := Order{}

	if query, args, errors := preloadInsert.Values(
		customer, generateStatus(customer.CreditCard),
	).ToSql(); errors != nil {
		return nil, errors
	} else if errors = database.Get(&order, query, args...); errors != nil {
		return nil, errors
	} else if customerCart, errors := insertProducts(database, order.ID, customerCartPOST); errors != nil {
		return nil, errors
	} else if token, errors := order.generateToken(); errors != nil {
		return nil, errors
	} else if updated, errors := updateOrder(database, token, customer.Username, map[string]interface{}{
		Token: token,
	}); errors != nil {
		return nil, errors
	} else {
		return updated.generateJson(customerCart), nil
	}
}

var preloadUpdate = common.SqlBuilder().Update(Orders).Suffix(common.ReturningRow)

func updateOrder(
	database *sqlx.DB,
	token string,
	customer string,
	what map[string]interface{},
) (*Order, error) {

	if query, args, errors := preloadUpdate.SetMap(what).Where(
		getQueryOptions(token, customer),
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var order Order
		return &order, database.Get(&order, query, args...)
	}
}

var preloadManyDelete = common.SqlBuilder().Delete(OrderProducts)
var preloadDelete = common.SqlBuilder().Delete(Orders).Suffix(common.ReturningRow)

func deleteOrder(database *sqlx.DB, token string, customer string) (sql.Result, error) {

	var order Order

	if query, args, errors := preloadDelete.Where(
		getQueryOptions(token, customer),
	).ToSql(); errors != nil {
		return nil, errors
	} else if errors := database.Get(&order, query, args...); errors != nil {
		return nil, errors
	} else {
		return preloadManyDelete.Where(squirrel.Eq{
			OrderID: order.ID,
		}).RunWith(database.DB).Exec()
	}
}

var preloadManyGet = common.SqlBuilder().Select(
	Quantity, "products.*",
).From(OrderProducts).Join(
	"products ON products.barcode = order_products.product_id",
)

func getProducts(database *sqlx.DB, orderId int) ([]CustomerCartJSON, error) {

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

func insertProducts(database *sqlx.DB, orderId int, customerCartPOST []CustomerCartPOST) ([]CustomerCartJSON, error) {

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