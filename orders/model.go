package orders

import (
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
	"github.com/marques999/acme-server/customers"
)

const (
	Orders        = "orders"
	Count         = "count"
	Token         = "token"
	Total         = "total"
	Customer      = "customer"
	Products      = "products"
	OrderProducts = "order_products"
	OrderID       = "order_id"
	ProductID     = "product_id"
	Quantity      = "quantity"
)

type Order struct {
	common.Model
	Customer string
	Token    string
}

type OrderJSON struct {
	common.Model
	Count    int     `binding:"required" json:"count"`
	Total    float64 `binding:"required" json:"total"`
	Customer string  `binding:"required" json:"customer"`
	Token    string  `binding:"required" json:"token"`
}

type OrderPOST struct {
	Signature string             `binding:"required" json:"signature"`
	Products  []CustomerCartPOST `binding:"required" json:"payload"`
}

type CustomerCartPOST struct {
	Product  string `binding:"required" json:"product"`
	Quantity int    `binding:"required" json:"quantity"`
}

type CustomerCartJSON struct {
	Quantity int                    `binding:"required" json:"quantity"`
	Product  products.ProductInsert `binding:"required" json:"product"`
}

func Migrate(database *sqlx.DB) {

	if _, errors := database.Exec(`CREATE TABLE orders(
		id SERIAL NOT NULL CONSTRAINT orders_pkey PRIMARY KEY,
		token TEXT DEFAULT FALSE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
		customer TEXT NOT NULL
			CONSTRAINT fk_orders_customer
			REFERENCES customers(username) ON UPDATE CASCADE ON DELETE CASCADE)
	`); errors != nil {
		return
	}

	if _, errors := database.Exec(`CREATE TABLE order_products(
		order_id INTEGER NOT NULL
			CONSTRAINT fk_order_products_order_id
			REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE,
		product_id INTEGER NOT NULL
			CONSTRAINT fk_order_products_product_id
			REFERENCES products(id) ON UPDATE CASCADE ON DELETE CASCADE,
		quantity INTEGER DEFAULT 1 NOT NULL,
		CONSTRAINT order_products_pkey PRIMARY KEY (order_id, product_id))
	`); errors != nil {
		return
	}

	creditCard := customers.CreditCard{
		Validity: time.Now().AddDate(5, 0, 0),
	}

	insertOrder(database, &customers.Customer{
		Username:   "admin",
		CreditCard: creditCard,
	}, []CustomerCartPOST{
		{"887899689185", 1},
	})

	insertOrder(database, &customers.Customer{
		Username:   "marques999",
		CreditCard: creditCard,
	}, []CustomerCartPOST{
		{"824142132142", 3},
		{"889349114872", 1},
	})

	insertOrder(database, &customers.Customer{
		Username:   "jabst",
		CreditCard: creditCard,
	}, []CustomerCartPOST{
		{"884102029028", 1},
		{"889349114872", 1},
	})

	insertOrder(database, &customers.Customer{
		Username:   "somouco",
		CreditCard: creditCard,
	}, []CustomerCartPOST{
		{"824142132142", 2},
	})
}