package creditcard

import (
	"time"
	"github.com/jmoiron/sqlx"
)

const (
	Type        = "type"
	Number      = "number"
	Validity    = "validity"
	CreditCards = "credit_cards"
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

func Migrate(database *sqlx.DB) {

	if _, errors := database.Exec(`CREATE TABLE credit_cards(
		id serial NOT NULL CONSTRAINT credit_cards_pkey PRIMARY KEY,
		type TEXT NOT NULL,
		number TEXT NOT NULL,
		validity timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL)
	`); errors == nil {

		Insert(database, &CreditCardJSON{
			Type:     "VISA",
			Number:   "123456789",
			Validity: time.Now().AddDate(5, 0, 0),
		})

		Insert(database, &CreditCardJSON{
			Type:     "Maestro",
			Number:   "310867542",
			Validity: time.Now().AddDate(3, 6, 0),
		})

		Insert(database, &CreditCardJSON{
			Type:     "Mastercard",
			Number:   "360420999",
			Validity: time.Now().AddDate(1, 3, 13),
		})

		Insert(database, &CreditCardJSON{
			Type:     "VISA Electron",
			Number:   "863101278",
			Validity: time.Now().AddDate(2, 5, 5),
		})
	}
}