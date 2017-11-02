package orders

import (
	"time"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"math/rand"
	"encoding/pem"
	"encoding/base64"
	"github.com/jmoiron/sqlx"
	"github.com/speps/go-hashids"
	"github.com/Masterminds/squirrel"
	"github.com/marques999/acme-server/common"
	"github.com/marques999/acme-server/products"
	"github.com/marques999/acme-server/customers"
)

func encodeSha1(payload []byte) []byte {

	sha1Algorithm := sha1.New()
	sha1Algorithm.Write([]byte(payload))

	return sha1Algorithm.Sum(nil)
}

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
		Status:           order.Status,
		Customer:         order.Customer,
		Products:         customerCart,
		common.CreatedAt: order.CreatedAt,
		common.UpdatedAt: order.UpdatedAt,
	}
}

func (order *OrderJSON) generateJson(
	customer *customers.Customer,
	customerCart []CustomerCartJSON,
) map[string]interface{} {

	return map[string]interface{}{
		Token:            order.Token,
		Count:            order.Count,
		Total:            order.Total,
		Status:           order.Status,
		Customer:         customer.GenerateDetails(&customer.CreditCard),
		Products:         customerCart,
		common.CreatedAt: order.CreatedAt,
		common.UpdatedAt: order.UpdatedAt,
	}
}

func generateStatus(creditCard customers.CreditCard) int {

	if creditCard.Validity.After(time.Now()) && rand.Float64() <= common.SuccessProbability {
		return ValidationComplete
	} else {
		return ValidationFailed
	}
}

func generateCustomerCart(query *sqlx.Rows) []CustomerCartJSON {

	var orderProducts []CustomerCartJSON

	for query.Next() {

		var quantity int
		var product products.Product

		query.Scan(&quantity, &product.ID, &product.Name,
			&product.Brand, &product.Price, &product.Barcode,
			&product.ImageUri, &product.Description,
			&product.CreatedAt, &product.UpdatedAt)

		orderProducts = append(orderProducts, CustomerCartJSON{
			quantity, product.GenerateJson(),
		})
	}

	return orderProducts
}

func verifySignature(publicKey string, signature string, payload []byte) error {

	pemCertificate := "-----BEGIN PUBLIC KEY-----\n" + publicKey + "\n-----END PUBLIC KEY-----"

	if decoded, errors := base64.StdEncoding.DecodeString(signature); errors != nil {
		return errors
	} else if publicKey, errors := decodePublicKey(pemCertificate); errors != nil {
		return errors
	} else {
		return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, encodeSha1(payload), decoded)
	}
}

func decodePublicKey(pemCertificate string) (*rsa.PublicKey, error) {

	if block, _ := pem.Decode([]byte(pemCertificate)); block != nil {

		publicKey, errors := x509.ParsePKIXPublicKey(block.Bytes)

		if errors == nil {

			switch publicKey := publicKey.(type) {
			case *rsa.PublicKey:
				return publicKey, nil
			}
		} else {
			return nil, errors
		}
	}

	return nil, common.PermissionDeniedError
}