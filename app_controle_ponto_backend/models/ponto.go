package models

import "time"

// Ponto representa um registro de ponto no sistema.
type Ponto struct {
	ID      string    `json:"id"`
	Horario time.Time `json:"horario"`
}
