package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	SourceBank          string          `validate:"required,min=2"`
	SourceAgency        string          `validate:"required,min=2"`
	SourceAccount       string          `validate:"required,min=2"`
	TargetBank          string          `validate:"required,min=2"`
	TargetAgency        string          `validate:"required,min=2"`
	TargetAccount       string          `validate:"required,min=2"`
	TransactionAmount   decimal.Decimal `validate:"required"`
	TransactionDatetime time.Time       `validate:"required"`
	CreatedAt           time.Time       `validate:"required"`
	UpdatedAt           time.Time       `validate:"required"`
}
