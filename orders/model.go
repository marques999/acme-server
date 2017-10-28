package orders

import (
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/creditcard"
	"github.com/marques999/acme-server/products"
)

const (
	ValidationFailed   = iota
	ValidationComplete = iota
	Purchased          = iota
	Orders             = "orders"
	Count              = "count"
	Token              = "token"
	Total              = "total"
	Status             = "status"
	Customer           = "customer"
	Products           = "products"
	OrderProducts      = "order_products"
	OrderID            = "order_products.order_id"
	ProductID          = "order_products.product_id"
	Quantity           = "order_products.quantity"
)

type OrderPOST struct {
	Signature string             `binding:"required" json:"signature"`
	Products  []CustomerCartPOST `binding:"required" json:"body"`
}

type CustomerCartPOST struct {
	Quantity int    `binding:"required" json:"quantity"`
	Product  string `binding:"required" json:"product"`
}

type CustomerCartJSON struct {
	Quantity int                  `binding:"required" json:"quantity"`
	Product  products.ProductJSON `binding:"required" json:"product"`
}

type Order struct {
	common.Model
	Status   int     `binding:"required" json:"status"`
	Count    int     `binding:"required" json:"count"`
	Total    float64 `binding:"required" json:"total"`
	Customer string  `binding:"required" json:"customer"`
	Token    string  `binding:"required" json:"token"`
}

func Migrate(database *sqlx.DB) {

	database.Exec(`CREATE TABLE orders(
		id serial NOT NULL CONSTRAINT orders_pkey PRIMARY KEY,
		created_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
		customer TEXT NOT NULL
			CONSTRAINT fk_orders_customer
			REFERENCES customers(username) ON UPDATE CASCADE ON DELETE CASCADE,
		total NUMERIC DEFAULT 0 NOT NULL,
		status INTEGER DEFAULT 0 NOT NULL,
		token TEXT DEFAULT FALSE NOT NULL)
	`)

	if _, errors := database.Exec(`CREATE TABLE order_products(
		order_id INTEGER NOT NULL
			CONSTRAINT fk_order_products_order_id
			REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE,
		product_id TEXT NOT NULL
			CONSTRAINT fk_order_products_product_id
			REFERENCES products(barcode) ON UPDATE CASCADE ON DELETE CASCADE,
		quantity INTEGER DEFAULT 1 NOT NULL,
		CONSTRAINT order_products_pkey PRIMARY KEY (order_id, product_id))
	`); errors != nil {
		return
	}

	creditCard := creditcard.CreditCard{
		Validity: time.Now().AddDate(5, 0, 0),
	}

	insertOrder(database, &customers.Customer{
		Username:   "admin",
		CreditCard: creditCard,
	}, []CustomerCartPOST{{
		Quantity: 1,
		Product:  "4713147489589",
	}})

	insertOrder(database, &customers.Customer{
		Username:   "marques999",
		CreditCard: creditCard,
	}, []CustomerCartPOST{
		{
			Quantity: 3,
			Product:  "824142132142",
		}, {
			Quantity: 1,
			Product:  "889349114872",
		},
	})

	insertOrder(database, &customers.Customer{
		Username:   "jabst",
		CreditCard: creditCard,
	}, []CustomerCartPOST{
		{
			Quantity: 1,
			Product:  "884102029028",
		}, {
			Quantity: 1,
			Product:  "889349114872",
		},
	})

	insertOrder(database, &customers.Customer{
		Username:   "somouco",
		CreditCard: creditCard,
	}, []CustomerCartPOST{{
		Quantity: 2,
		Product:  "824142132142",
	}})
}