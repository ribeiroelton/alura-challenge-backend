package spi

import (
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
)

type TransactionRepository interface {
	SaveTransaction(*model.Transaction) error
	ListTransactions() ([]model.Transaction, error)
}
