package orders

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	"github.com/speps/go-hashids"
	"github.com/marques999/acme-server/common"
)

func encodeSha1(payload []byte) []byte {
	sha1Algorithm := sha1.New()
	sha1Algorithm.Write([]byte(payload))
	return sha1Algorithm.Sum(nil)
}

func getQueryOptions(orderId string, customer string) map[string]interface{} {

	if customer == common.AdminAccount {
		return map[string]interface{}{
			"id": orderId,
		}
	} else {
		return map[string]interface{}{
			"id":       orderId,
			"customer": customer,
		}
	}
}

func decodePublicKey(publicKey string) (key *rsa.PublicKey) {

	block, _ := pem.Decode([]byte(publicKey))

	if block != nil {

		if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {

			switch pub := pub.(type) {
			case *rsa.PublicKey:
				return pub
			}
		}
	}

	return nil
}

func GenerateToken(order *Order) (string, error) {

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

func generateJson(order Order) map[string]interface{} {

	return gin.H{
		"id":       order.ID,
		"token":    order.Token,
		"total":    order.Total,
		"status":   order.Status,
		"created":  order.CreatedAt,
		"modified": order.UpdatedAt,
		"products": order.Products,
	}
}