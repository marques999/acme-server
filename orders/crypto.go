package orders

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/marques999/acme-server/common"
)

func encodeSha1(payload []byte) []byte {

	sha1Algorithm := sha1.New()
	sha1Algorithm.Write([]byte(payload))

	return sha1Algorithm.Sum(nil)
}

func verifySignature(publicKey string, signature string, payload []byte) error {

	if decoded, errors := base64.StdEncoding.DecodeString(signature); errors != nil {
		return errors
	} else if publicKey, errors := decodePublicKey(publicKey); errors != nil {
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