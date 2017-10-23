package common

import "time"

type Model struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"modified"`
	DeletedAt *time.Time `json:"deleted" gorm:"index"`
}

type Encrypted struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}
