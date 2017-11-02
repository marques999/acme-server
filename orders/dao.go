package orders

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
	"github.com/marques999/acme-server/customers"
	"time"
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
	"products ON products.id = order_products.product_id",
).GroupBy("orders.id")

func listOrders(database *sqlx.DB, username string) ([]OrderJSON, error) {

	builder := preloadGet

	if username != common.AdminAccount {
		builder = builder.Where(squirrel.Eq{Customer: username})
	}

	if query, args, errors := builder.ToSql(); errors != nil {
		return []OrderJSON{}, nil
	} else {
		var orders []OrderJSON
		return orders, database.Select(&orders, query, args...)
	}
}

func getOrder(database *sqlx.DB, token string, customer string) (*OrderJSON, error) {

	if query, args, errors := preloadGet.Where(getQueryOptions(
		token, customer,
	)).ToSql(); errors != nil {
		return nil, errors
	} else {
		var order OrderJSON
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
		customer.Username, generateStatus(customer.CreditCard),
	).ToSql(); errors != nil {
		println("InsertProducts1")
		return nil, errors
	} else if errors = database.Get(&order, query, args...); errors != nil {
		println(errors.Error())
		return nil, errors
	} else if customerCart, errors := insertProducts(database, order.ID, customerCartPOST); errors != nil {
		println("InsertProducts3")
		return nil, errors
	} else if token, errors := order.generateToken(); errors != nil {
		println("InsertProducts4")
		return nil, errors
	} else if query, args, errors := preloadUpdate.Set(Token, token).Where(
		squirrel.Eq{common.Id: order.ID},
	).ToSql(); errors != nil {
		println("InsertProduct5")
		return nil, errors
	} else if errors := database.Get(&order, query, args...); errors != nil {
		println("InsertProducts6")
		return nil, errors
	} else {
		return order.generateJson(customerCart), nil
	}
}

var preloadManyDelete = common.SqlBuilder().Delete(OrderProducts)
var preloadDelete = common.SqlBuilder().Delete(Orders).Suffix(common.ReturningRow)
var preloadUpdate = common.SqlBuilder().Update(Orders).Suffix(common.ReturningRow)

func updateOrder(database *sqlx.DB, token string) (*Order, error) {

	if query, args, errors := preloadUpdate.SetMap(map[string]interface{}{
		Status:           Purchased,
		common.UpdatedAt: time.Now(),
	}).Where(
		squirrel.Eq{Token: token},
	).ToSql(); errors != nil {
		return nil, errors
	} else {
		var order Order
		return &order, database.Get(&order, query, args...)
	}
}

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
	"products ON products.id = order_products.product_id",
)

func getProducts(database *sqlx.DB, orderId int) ([]CustomerCartJSON, error) {

	if query, args, errors := preloadManyGet.Where(
		squirrel.Eq{OrderID: orderId},
	).ToSql(); errors != nil {
		return []CustomerCartJSON{}, errors
	} else if items, errors := database.Queryx(query, args...); errors != nil {
		return []CustomerCartJSON{}, errors
	} else {
		return generateCustomerCart(items), nil
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
	purchased, errors := products.GetProductsByBarcode(database, barcodes)

	if errors != nil {
		return customerCart, errors
	}

	builder := preloadManyInsert

	for index, product := range purchased {

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