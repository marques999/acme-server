package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/speps/go-hashids"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
)

func getQueryOptions(orderId string, customer string) squirrel.Eq {

	if customer == common.AdminAccount {
		return squirrel.Eq{
			Token: orderId,
		}
	} else {
		return squirrel.Eq{
			Token:      orderId,
			"customer": customer,
		}
	}
}

func (order *Order) generateToken() (string, error) {

	hashData := hashids.NewData()
	hashData.MinLength = 8
	hashData.Salt = common.RamenRecipe
	hashGenerator, _ := hashids.NewWithData(hashData)

	return hashGenerator.Encode([]int{
		order.ID,
		order.CreatedAt.Hour(),
		order.CreatedAt.Minute(),
	})
}

func (order *Order) generateList() *map[string]interface{} {

	return &map[string]interface{}{
		Status:           order.Status,
		Count:            order.Count,
		Total:            order.Total,
		Customer:         order.Customer,
		Token:            order.Token,
		common.CreatedAt: order.CreatedAt,
		common.UpdatedAt: order.UpdatedAt,
	}
}

func (order *Order) generateJson(customerCart []CustomerCartJSON) *map[string]interface{} {

	return &map[string]interface{}{
		Token:            order.Token,
		Total:            order.Total,
		Status:           order.Status,
		Customer:         order.Customer,
		Products:         customerCart,
		common.CreatedAt: order.CreatedAt,
		common.UpdatedAt: order.UpdatedAt,
	}
}

func generateProductOrder(query *sqlx.Rows) CustomerCartJSON {

	var quantity int
	var product products.Product

	query.Scan(&quantity)
	query.StructScan(&product)

	return CustomerCartJSON{
		Quantity: quantity,
		Product:  product.GenerateJson(),
	}
}

func generateCustomerCart(query *sqlx.Rows) []CustomerCartJSON {

	orderProducts := []CustomerCartJSON{}

	for query.Next() {
		orderProducts = append(orderProducts, generateProductOrder(query))
	}

	return orderProducts
}