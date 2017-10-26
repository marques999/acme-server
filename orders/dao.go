package orders

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
)

func listOrders(db *sqlx.DB, username string) ([]Order, error) {

	/*	sqlQuery, sqlArgs, sqlException := squirrel.Select(Orders).Columns("*").Where(
			squirrel.Eq{CustomerID: username},
		).ToSql()

		if sqlException != nil {

		}*/
	orders := []Order{}
	return orders, nil
}

func deleteOrder(db *sqlx.DB, token string) (sql.Result, error) {
	return squirrel.Delete(Orders).Where(squirrel.Eq{Token: token}).RunWith(db).Exec()
}

func insertOrder(database *sqlx.DB, customer string, products ...CustomerCartPOST) (*OrderJSON, error) {

	// inserir encomenda sem token nem total

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Insert(Orders).Columns(
		Customer, Status, Total,
	).Values(
		customer, ValidationFailed, 0.0,
	).Suffix("RETURNING *").ToSql()

	if sqlException != nil {
		return nil, sqlException
	}

	// obtém encomenda recentemente inserida

	order := Order{}
	sqlException = database.Get(&order, sqlQuery, sqlArgs...)

	if sqlException != nil {
		return nil, sqlException
	}

	// inserir as entidades da associação muitos-para-muitos

	customerCart, sqlException := insertOrderProducts(database, order.ID, products)

	if sqlException != nil {
		return nil, sqlException
	}

	// atualizar encomenda com total a pagar e token gerado aleatoriamente

	orderToken, hashException := order.generateToken()

	if hashException != nil {
		return nil, sqlException
	}

	sqlQuery, sqlArgs, sqlException = common.StatementBuilder().Update(Orders).Set(
		Total, calculateTotal(customerCart)).Set(Token, orderToken,
	).Where(squirrel.Eq{common.Id: order.ID}).Suffix("RETURNING *").ToSql()

	if sqlException != nil {
		return nil, sqlException
	}

	// obter encomenda atualizada, retornando a função

	sqlException = database.Get(&order, sqlQuery, sqlArgs...)

	if sqlException != nil {
		return nil, sqlException
	} else {
		return order.generateJson(customerCart), nil
	}
}

func removeOrderProducts(db *sqlx.DB, orderId int) (sql.Result, error) {
	return squirrel.Delete(OrderProducts).Where(squirrel.Eq{OrderID: orderId}).RunWith(db).Exec()
}

func insertOrderProducts(db *sqlx.DB, orderId int, productOrderPOST []CustomerCartPOST) ([]CustomerCartJSON, error) {

	barcodes := make([]string, len(productOrderPOST))

	for index, productOrder := range productOrderPOST {
		barcodes[index] = productOrder.Product
	}

	productOrders := make([]CustomerCartJSON, len(productOrderPOST))
	customerCart, sqlException := products.GetProductsByBarcode(db, barcodes)

	if sqlException != nil {
		return productOrders, sqlException
	}

	builder := squirrel.Insert(OrderProducts).Columns(OrderID, ProductID, Quantity)

	for index, product := range customerCart {

		quantity := productOrderPOST[index].Quantity
		builder = builder.Values(orderId, product.ID, quantity)

		productOrders[index] = CustomerCartJSON{
			Product:  product.GenerateJson(),
			Quantity: quantity,
		}
	}

	_, sqlException = builder.RunWith(db).Exec()

	if sqlException != nil {
		return []CustomerCartJSON{}, sqlException
	} else {
		return productOrders, nil
	}
}

func getOrderProducts(database *sqlx.DB, orderId int) ([]CustomerCartJSON, error) {

	sqlQuery, sqlArgs, sqlException := common.StatementBuilder().Select(
		"products.*",
	).From(OrderProducts).Join(
		"products USING(product_id)",
	).Where(squirrel.Eq{OrderID: orderId}).ToSql()

	if sqlException != nil {
		return []CustomerCartJSON{}, sqlException
	}

	orderProducts := []CustomerCartJSON{}
	sqlRows, sqlException := database.Queryx(sqlQuery, sqlArgs)

	if sqlException != nil {
		return orderProducts, sqlException
	}

	for sqlRows.Next() {

		var quantity int
		var product products.Product

		if sqlException := sqlRows.Scan(&quantity); sqlException != nil {
			return orderProducts, sqlException
		} else if sqlException := sqlRows.StructScan(&product); sqlException != nil {
			return orderProducts, sqlException
		} else {
			orderProducts = append(orderProducts, CustomerCartJSON{
				Quantity: quantity,
				Product:  product.GenerateJson(),
			})
		}
	}

	return orderProducts, nil
}