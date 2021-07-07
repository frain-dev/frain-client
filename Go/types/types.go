package types

import "time"

type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    []Component `json:"data"`
}

type Component struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
