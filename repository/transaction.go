package repository

import (
	"context"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type TransactionRepository struct {
	db *mongo.Database
}

const (
	transactionCollection = "transactions"
)

func NewTransactionRepository(c *config.Config) (*TransactionRepository, error) {
	opts := options.Client().ApplyURI(c.ConnString)
	cli, err := mongo.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cli.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err := cli.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	db := cli.Database(c.DatabaseName)

	return &TransactionRepository{
		db: db,
	}, nil
}

//SaveTransaction saves a new transaction to the database. Default timeout of 3 seconds
func (tr *TransactionRepository) SaveTransaction(t *model.Transaction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tr.db.Collection(transactionCollection).InsertOne(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

//ListTransactions returns a slice of Transactions. Empty if no transactions found. Default timeout of 30 seconds
func (tr *TransactionRepository) ListTransactions() ([]model.Transaction, error) {
	ts := []model.Transaction{}

	return ts, nil
}
