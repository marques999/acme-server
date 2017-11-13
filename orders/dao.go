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
	"orders.token",
).From(Orders).Join(
	"order_products ON orders.id = order_products.order_id",
).Join(
	"products ON products.id = order_products.product_id",
).GroupBy(
	"orders.id",
).OrderBy(
	"orders.created_at DESC",
)

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

var preloadUpdate = common.SqlBuilder().Update(Orders)
var preloadInsert = common.SqlBuilder().Insert(Orders).Columns(Customer).Suffix(common.ReturningRow)

func insertOrder(
	database *sqlx.DB,
	customer *customers.Customer,
	customerCartPOST []CustomerCartPOST,
) (*OrderJSON, error) {

	order := Order{}

	if verifyPurchase(customer.CreditCard) == false {
		return nil, common.PurchaseValidationError
	} else if query, args, errors := preloadInsert.Values(customer.Username).ToSql(); errors != nil {
		return nil, errors
	} else if errors = database.Get(&order, query, args...); errors != nil {
		return nil, errors
	} else if errors := insertProducts(database, order.ID, customerCartPOST); errors != nil {
		return nil, errors
	} else if token, errors := order.generateToken(); errors != nil {
		return nil, errors
	} else if _, errors := preloadUpdate.Set(Token, token).Where(
		squirrel.Eq{common.Id: order.ID},
	).RunWith(database.DB).Exec(); errors != nil {
		return nil, errors
	} else if order, errors := getOrder(database, token, customer.Username); errors != nil {
		return nil, errors
	} else {
		return order, nil
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

func insertProducts(database *sqlx.DB, orderId int, customerCartPOST []CustomerCartPOST) error {

	builder := preloadManyInsert
	barcodes := make([]string, len(customerCartPOST))

	for index, entry := range customerCartPOST {
		barcodes[index] = entry.Product
	}

	purchased, errors := products.GetProductsByBarcode(database, barcodes)

	if errors != nil {
		return errors
	}

	for index, product := range purchased {
		builder = builder.Values(orderId, product.ID, customerCartPOST[index].Quantity)
	}

	if _, errors = builder.RunWith(database.DB).Exec(); errors != nil {
		return errors
	} else {
		return nil
	}
}