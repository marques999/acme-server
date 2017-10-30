package customers

import (
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/common"
)

const (
	Type           = "type"
	Number         = "number"
	Validity       = "validity"
	Customers      = "customers"
	Name           = "name"
	Country        = "country"
	Username       = "username"
	Password       = "password"
	Address1       = "address1"
	Address2       = "address2"
	PublicKey      = "public_key"
	TaxNumber      = "tax_number"
	CreditCards    = "credit_cards"
	CreditCardData = "credit_card"
	CreditCardID   = "credit_card_id"
)

type CreditCard struct {
	ID       int
	Type     string
	Number   string
	Validity time.Time
}

type CreditCardJSON struct {
	Type     string    `binding:"required" json:"type"`
	Number   string    `binding:"required" json:"number"`
	Validity time.Time `binding:"required" json:"validity"`
}

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
	CreditCard   CreditCard
}

type CustomerList struct {
	Name      string    `binding:"required" json:"name"`
	Country   string    `binding:"required" json:"country"`
	Username  string    `binding:"required" json:"username"`
	Address1  string    `binding:"required" json:"address1"`
	Address2  string    `binding:"required" json:"address2"`
	TaxNumber string    `binding:"required" json:"tax_number" db:"tax_number"`
	CreatedAt time.Time `binding:"required" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `binding:"required" json:"updated_at" db:"updated_at"`
}

type CustomerInsert struct {
	Name       string         `binding:"required" json:"name"`
	Country    string         `binding:"required" json:"country"`
	Username   string         `binding:"required" json:"username"`
	Password   string         `binding:"required" json:"password"`
	Address1   string         `binding:"required" json:"address1"`
	Address2   string         `binding:"required" json:"address2"`
	PublicKey  string         `binding:"required" json:"public_key"`
	TaxNumber  string         `binding:"required" json:"tax_number"`
	CreditCard CreditCardJSON `binding:"required" json:"credit_card"`
}

type CustomerUpdate struct {
	Name       string         `binding:"required" json:"name"`
	Country    string         `binding:"required" json:"country"`
	Password   string         `binding:"required" json:"password"`
	Address1   string         `binding:"required" json:"address1"`
	Address2   string         `binding:"required" json:"address2"`
	PublicKey  string         `binding:"required" json:"public_key"`
	TaxNumber  string         `binding:"required" json:"tax_number"`
	CreditCard CreditCardJSON `binding:"required" json:"credit_card"`
}

func Migrate(database *sqlx.DB) {

	if _, errors := database.Exec(`CREATE TABLE credit_cards(
		id serial NOT NULL CONSTRAINT credit_cards_pkey PRIMARY KEY,
		type TEXT NOT NULL,
		number TEXT NOT NULL,
		validity timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL)
	`); errors == nil {

		insertCreditCard(database, &CreditCardJSON{
			"VISA", "123456789",
			time.Now().AddDate(5, 0, 0),
		})

		insertCreditCard(database, &CreditCardJSON{
			"Maestro", "310867542",
			time.Now().AddDate(3, 6, 0),
		})

		insertCreditCard(database, &CreditCardJSON{
			"Mastercard", "360420999",
			time.Now().AddDate(1, 3, 13),
		})

		insertCreditCard(database, &CreditCardJSON{
			"VISA Electron", "863101278",
			time.Now().AddDate(2, 5, 5),
		})
	}

	if _, errors := database.Exec(`CREATE TABLE customers(
		id serial NOT NULL CONSTRAINT customers_pkey PRIMARY KEY,
		name TEXT NOT NULL,
		country VARCHAR(2) NOT NULL,
		username VARCHAR(32) NOT NULL UNIQUE,
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
	`); errors == nil {

		database.MustExec("CREATE INDEX IF NOT EXISTS idx_customers_username ON customers(username)")
		user1Password, _ := common.GeneratePassword("admin")
		user2Password, _ := common.GeneratePassword("r0wsauce")
		user3Password, _ := common.GeneratePassword("bighotshaq")
		user4Password, _ := common.GeneratePassword("skibidipoop")

		insertCustomer(database, CustomerInsert{
			Name:      "Administrator",
			Username:  "admin",
			Password:  user1Password,
			TaxNumber: "930248516",
			Address1:  "Rua Branco, Nº 25",
			Address2:  "8681-962 Tomar",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAL1L9h1N9xqNe0I4ddyjKD6lv0ArcEhBJbU550urvmvJ
qa1Rm8Zr+V0+VCp9swcCAwEAAQ==`,
		}, 1)

		insertCustomer(database, CustomerInsert{
			Name:      "Diogo Marques",
			Username:  "marques999",
			Password:  user2Password,
			TaxNumber: "761489053",
			Address1:  "Rua São Diogo, Nº 855",
			Address2:  "6311-969 Vendas Novas",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAKCRuhMUuFoJvDVeicvyfyQf9ADQ1qNe+dabNSpOkr76
FcVTBd+TBe2sEshVefUCAwEAAQ==`,
		}, 2)

		insertCustomer(database, CustomerInsert{
			Username:  "jabst",
			Password:  user3Password,
			Name:      "José Teixeira",
			TaxNumber: "685102439",
			Address1:  "Avenida Lima, Nº 167",
			Address2:  "7049-952 Santa Cruz",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvALLIEFJe1v3hiGpzYlzo/hxEXBW2XrA47b/S2i0X7ZZv
08HLhNfdPr2XC8ZzLpECAwEAAQ==`,
		}, 3)

		insertCustomer(database, CustomerInsert{
			Username:  "somouco",
			Name:      "Carlos Samouco",
			Password:  user4Password,
			TaxNumber: "537812640",
			Address1:  "Travessa Mia Assunção, Nº 532",
			Address2:  "5334-964 Coimbra",
			Country:   "PT",
			PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAK0smd9hF2yMJOeidEDq2GieQJY2Ac3bRpoXeOpiD/Oi
pBrNyqlMpzEKUF917T0CAwEAAQ==`,
		}, 4)
	}
}