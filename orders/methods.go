package orders

import (
	"crypto/sha1"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/speps/go-hashids"
)

func encodeSha1(payload []byte) []byte {
	sha1Algorithm := sha1.New()
	sha1Algorithm.Write([]byte(payload))
	return sha1Algorithm.Sum(nil)
}

func getQueryOptions(orderId string, customerId string) map[string]interface{} {

	if customerId == "admin" {
		return map[string]interface{}{
			"id": orderId,
		}
	} else {
		return map[string]interface{}{
			"id":       orderId,
			"customer": customerId,
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

func GenerateHashId(order *Order) (string, error) {

	hashIds := hashids.NewData()
	hashIds.Salt = "acmestore"
	hashIds.MinLength = 6
	hashGenerator, hashException := hashids.NewWithData(hashIds)

	if hashException != nil {
		return "", hashException
	} else {
		return hashGenerator.Encode([]int{
			order.ID,
			order.CreatedAt.Hour(),
			order.CreatedAt.Minute(),
		})
	}
}
