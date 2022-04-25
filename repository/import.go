package repository

import (
	"context"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ImportRepository struct {
	db *mongo.Database
}

const (
	importsCollection = "imports"
)

func NewImportRepository(c *config.Config) (*ImportRepository, error) {
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

	return &ImportRepository{
		db: db,
	}, nil
}

//SaveImport saves a import stats to the database. Default timeout of 3 seconds
func (tr *ImportRepository) SaveImport(u *model.Import) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tr.db.Collection(importsCollection).InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

//HasImportByTransactionDate returns true if has transactions imports in a specific date. Default timeout of 3 seconds.
func (tr *ImportRepository) HasImportByTransactionDate(dt time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	date, _ := time.Parse("2006-01-02", dt.Format("2006-01-02"))

	filter := bson.D{{Key: "transactiondate", Value: bson.D{{Key: "$gte", Value: date}, {Key: "$lt", Value: date.AddDate(0, 0, 1)}}}}

	t := tr.db.Collection(importsCollection).FindOne(ctx, filter)
	if t.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if t.Err() != nil {
		return false, t.Err()
	}

	return true, nil
}

//ListImports returns a slice of imports. Default timeout of 30 seconds.
func (tr *ImportRepository) ListImports() ([]model.Import, error) {
	us := []model.Import{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := tr.db.Collection(importsCollection).Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		u := model.Import{}
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	return us, nil
}
