package common

import "time"

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"modified"`
	DeletedAt *time.Time `json:"deleted" gorm:"index"`
}
