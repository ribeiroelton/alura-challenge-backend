package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                  string          `bson:"_id,omitempty"`
	SourceBank          string          `validate:"required,min=2" bson:"source_bank"`
	SourceAgency        string          `validate:"required,min=2" bson:"source_agency"`
	SourceAccount       string          `validate:"required,min=2" bson:"source_account"`
	TargetBank          string          `validate:"required,min=2" bson:"target_bank"`
	TargetAgency        string          `validate:"required,min=2" bson:"target_agency"`
	TargetAccount       string          `validate:"required,min=2" bson:"target_account"`
	TransactionAmount   decimal.Decimal `validate:"required" bson:"transaction_amount"`
	TransactionDatetime time.Time       `validate:"required" bson:"transaction_date_time"`
	CreatedAt           time.Time       `validate:"required" bson:"created_at"`
	UpdatedAt           time.Time       `validate:"required" bson:"updated_at"`
}

type Import struct {
	ID              string    `bson:"_id,omitempty"`
	TransactionDate time.Time `validate:"required" bson:"transaction_date"`
	ImportDate      time.Time `validate:"required" bson:"import_date"`
	Filename        string    `validate:"required,min=1" bson:"filename"`
	FileSizeInMB    float64   `validate:"required" bson:"file_size"`
	CreatedAt       time.Time `validate:"required" bson:"created_at"`
	UpdatedAt       time.Time `validate:"required" bson:"updated_at"`
}
