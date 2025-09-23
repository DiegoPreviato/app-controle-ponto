package models

import "time"

// Ponto representa um registro de ponto no sistema.
type Ponto struct {
	ID      string    `json:"id"`
	UserID  int64     `json:"user_id"`
	Horario time.Time `json:"horario"`
}
