package common

import "time"

type Model struct {
	ID        int       `sql:"primary_key"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp;not null"`
	UpdatedAt time.Time `sql:"DEFAULT:current_timestamp;not null"`
}