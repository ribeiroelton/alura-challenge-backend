package spi

import (
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
)

type TransactionRepository interface {
	SaveTransaction(*model.Transaction) error
	ListTransactions() ([]model.Transaction, error)
	SaveImport(*model.Import) error
	HasImportByTransactionDate(time.Time) (bool, error)
	ListImports() ([]model.Import, error)
}
