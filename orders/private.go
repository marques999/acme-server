package orders

import (
	"time"
	"math/rand"
	"github.com/jmoiron/sqlx"
	"github.com/speps/go-hashids"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/creditcard"
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

func generateStatus(creditCard creditcard.CreditCard) int {

	if creditCard.Validity.After(time.Now()) && rand.Float64() <= common.SuccessProbability {
		return ValidationComplete
	} else {
		return ValidationFailed
	}
}

func generateCustomerCart(query *sqlx.Rows) []CustomerCartJSON {

	orderProducts := []CustomerCartJSON{}

	for query.Next() {

		var quantity int
		var product products.Product

		query.Scan(&quantity)
		query.StructScan(&product)

		orderProducts = append(orderProducts, CustomerCartJSON{
			Quantity: quantity,
			Product:  product.GenerateJson(),
		})
	}

	return orderProducts
}