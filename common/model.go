package common

import "time"

const (
	Id                  = "id"
	CreatedAt           = "created_at"
	UpdatedAt           = "updated_at"
	RouteDefault        = "/"
	RouteWithId         = "/:id"
	ReturningRow        = "RETURNING *"
	RamenRecipe         = "mieic@feup#2017"
	AdminAccount        = "admin"
	SuccessProbability  = 0.95
	AuthenticationRealm = "fe.up.pt"
)

type Model struct {
	ID        int       `binding:"required" json:"id"`
	CreatedAt time.Time `binding:"required" db:"created_at" json:"created_at"`
	UpdatedAt time.Time `binding:"required" db:"updated_at" json:"updated_at"`
}