package customers

import (
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/creditcard"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/common"
	"time"
)

const (
	Customers  = "customers"
	Name       = "name"
	Country    = "country"
	Username   = "username"
	Password   = "password"
	Address1   = "address1"
	Address2   = "address2"
	PublicKey  = "public_key"
	TaxNumber  = "tax_number"
	CreditCard = "credit_card_id"
)

type Customer struct {
	common.Model
	Name         string
	Country      string
	Username     string
	Password     string
	Address1     string
	Address2     string
	PublicKey    string `db:"public_key"`
	TaxNumber    string `db:"tax_number"`
	CreditCardID int    `db:"credit_card_id"`
}

type CustomerJSON struct {
	Name       string                    `binding:"required" json:"name"`
	Country    string                    `binding:"required" json:"country"`
	Username   string                    `binding:"required" json:"username"`
	Address1   string                    `binding:"required" json:"address1"`
	Address2   string                    `binding:"required" json:"address2"`
	TaxNumber  string                    `binding:"required" json:"tax_number"`
	CreditCard *creditcard.CreditCardJSON `binding:"required" json:"credit_card"`
	CreatedAt  time.Time                 `binding:"required" json:"created_at"`
	UpdatedAt  time.Time                 `binding:"required" json:"updated_at"`
}

type CustomerPOST struct {
	Name       string                    `binding:"required" json:"name"`
	Country    string                    `binding:"required" json:"country"`
	Username   string                    `binding:"required" json:"username"`
	Password   string                    `binding:"required" json:"password"`
	Address1   string                    `binding:"required" json:"address1"`
	Address2   string                    `binding:"required" json:"address2"`
	PublicKey  string                    `binding:"required" json:"public_key"`
	TaxNumber  string                    `binding:"required" json:"tax_number"`
	CreditCard creditcard.CreditCardJSON `binding:"required" json:"credit_card"`
}

func Migrate(database *sqlx.DB) {

	if _, sqlException := database.Exec(`CREATE TABLE customers(
		id serial NOT NULL CONSTRAINT customers_pkey PRIMARY KEY,
		name TEXT NOT NULL,
		country VARCHAR(2) NOT NULL,
		username VARCHAR(32) NOT NULL CONSTRAINT customers_username_key UNIQUE,
		password VARCHAR(64) NOT NULL,
		address1 TEXT NOT NULL,
		address2 TEXT NOT NULL,
		public_key TEXT NOT NULL,
		tax_number VARCHAR(9),
		created_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
		credit_card_id INTEGER
			CONSTRAINT fk_customers_credit_card_id
			REFERENCES credit_cards(id) ON UPDATE CASCADE ON DELETE CASCADE)
	`); sqlException == nil {

		insertCustomer(database, CustomerPOST{
			Name:      "Administrator",
			Username:  "admin",
			Password:  auth.KamikazePassword("admin"),
			TaxNumber: "930248516",
			Address1:  "Rua Branco, Nº 25",
			Address2:  "8681-962 Tomar",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAL1L9h1N9xqNe0I4ddyjKD6lv0ArcEhBJbU550urvmvJ
qa1Rm8Zr+V0+VCp9swcCAwEAAQ==`,
		}, 1)

		insertCustomer(database, CustomerPOST{
			Name:      "Diogo Marques",
			Username:  "marques999",
			Password:  auth.KamikazePassword("r0wsauce"),
			TaxNumber: "761489053",
			Address1:  "Rua São Diogo, Nº 855",
			Address2:  "6311-969 Vendas Novas",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAKCRuhMUuFoJvDVeicvyfyQf9ADQ1qNe+dabNSpOkr76
FcVTBd+TBe2sEshVefUCAwEAAQ==`,
		}, 2)

		insertCustomer(database, CustomerPOST{
			Username:  "jabst",
			Password:  auth.KamikazePassword("bighotshaq"),
			Name:      "José Teixeira",
			TaxNumber: "685102439",
			Address1:  "Avenida Lima, Nº 167",
			Address2:  "7049-952 Santa Cruz",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvALLIEFJe1v3hiGpzYlzo/hxEXBW2XrA47b/S2i0X7ZZv
08HLhNfdPr2XC8ZzLpECAwEAAQ==`,
		}, 3)

		insertCustomer(database, CustomerPOST{
			Username:  "somouco",
			Name:      "Carlos Samouco",
			Password:  auth.KamikazePassword("skibidipap"),
			TaxNumber: "537812640",
			Address1:  "Travessa Mia Assunção, Nº 532",
			Address2:  "5334-964 Coimbra",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAK0smd9hF2yMJOeidEDq2GieQJY2Ac3bRpoXeOpiD/Oi
pBrNyqlMpzEKUF917T0CAwEAAQ==`,
		}, 4)
	}
}
