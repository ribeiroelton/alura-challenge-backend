package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//MongoDB struct that implements all Repository spi interface
type MongoDB struct {
	DB  *mongo.Database
	CLI *mongo.Client
}

//NewMongoDB creates a new MongoDB
func NewMongoDB(c *config.Config) (*MongoDB, error) {
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

	return &MongoDB{DB: db, CLI: cli}, nil
}

//SaveTransaction saves a new transaction to the database. Default timeout of 3 seconds
func (m *MongoDB) SaveTransaction(t *model.Transaction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Collection("transactions").InsertOne(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

//ListTransactions returns a slice of Transactions. Empty if no transactions found. Default timeout of 30 seconds
func (m *MongoDB) ListTransactions() ([]model.Transaction, error) {
	ts := []model.Transaction{}

	return ts, nil
}

//SaveImport saves a import stats to the database. Default timeout of 3 seconds
func (m *MongoDB) SaveImport(u *model.Import) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Collection("imports").InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

//HasImportByTransactionDate returns true if has transactions imports in a specific date. Default timeout of 3 seconds.
func (m *MongoDB) HasImportByTransactionDate(dt time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	date, _ := time.Parse("2006-01-02", dt.Format("2006-01-02"))

	filter := bson.D{{Key: "transactiondate", Value: bson.D{{Key: "$gte", Value: date}, {Key: "$lt", Value: date.AddDate(0, 0, 1)}}}}

	r := m.DB.Collection("imports").FindOne(ctx, filter)
	if r.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if r.Err() != nil {
		return false, r.Err()
	}

	return true, nil
}

//ListImports returns a slice of imports. Default timeout of 30 seconds.
func (m *MongoDB) ListImports() ([]model.Import, error) {
	us := []model.Import{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := m.DB.Collection("imports").Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		u := model.Import{}
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		us = append(us, u)
		log.Println(u)
	}

	return us, nil
}

//SaveUser saves a new user to the database if the specified email does not already exists.
func (m *MongoDB) SaveUser(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

//UpdateUserByEmail updates an user using its email as key. Default timeout of 3 seconds
//TODO implements DeleteUserByEmail
func (m *MongoDB) UpdateUserByEmail(update *model.User) error {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return nil
}

//DeleteUserByEmail deletes an user by using its email. Default timeout of 3 seconds
func (m *MongoDB) DeleteUserByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: email, Value: email}}

	r, err := m.DB.Collection("users").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if r.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

//GetUserByEmail returns an user by its email. Default timeout of 3 seconds
func (m *MongoDB) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: email, Value: email}}

	r := m.DB.Collection("users").FindOne(ctx, filter)
	if r.Err() != nil {
		return nil, r.Err()
	}

	u := &model.User{}

	if err := r.Decode(u); err != nil {
		return nil, err
	}

	return u, nil
}

//HasUserByEmail returns true if found and user by its email, otherwise returns false. Default timeout of 3 seconds
func (m *MongoDB) HasUserByEmail(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: email, Value: email}}

	r := m.DB.Collection("users").FindOne(ctx, filter)
	if r.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if r.Err() != nil {
		return false, r.Err()
	}

	return true, nil
}

//ListUsers returns a slice of users. Empty if no users found. Default timeout of 30 seconds.
func (m *MongoDB) ListUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	users := []model.User{}

	cur, err := m.DB.Collection("users").Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	ctxList, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for cur.Next(ctxList) {
		u := &model.User{}

		if err := cur.Decode(u); err != nil {
			return nil, err
		}

		users = append(users, *u)
	}

	return users, nil
}
