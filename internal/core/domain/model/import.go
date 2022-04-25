package model

import "time"

type Import struct {
	TransactionDate time.Time `validate:"required"`
	ImportDate      time.Time `validate:"required"`
	Filename        string    `validate:"required,min=1"`
	FileSizeInMB    float64   `validate:"required"`
	CreatedAt       time.Time `validate:"required"`
	UpdatedAt       time.Time `validate:"required"`
}
