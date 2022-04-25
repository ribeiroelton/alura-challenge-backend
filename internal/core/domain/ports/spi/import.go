package spi

import (
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
)

type ImportRepository interface {
	SaveImport(*model.Import) error
	HasImportByTransactionDate(time.Time) (bool, error)
	ListImports() ([]model.Import, error)
}
