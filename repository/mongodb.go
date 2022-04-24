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

type MongoDB struct {
	DB  *mongo.Database
	CLI *mongo.Client
}

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

func (m *MongoDB) SaveTransaction(t *model.Transaction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Collection("transactions").InsertOne(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) ListTransactionsByDate() ([]model.Transaction, error) {
	ts := []model.Transaction{}

	return ts, nil
}

func (m *MongoDB) SaveImport(u *model.Import) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ok, err := m.HasImportByTransactionDate(u.TransactionDate)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	_, err = m.DB.Collection("imports").InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

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

func (m *MongoDB) SaveUser(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

// func (m *MongoDB) UpdateUserByEmail(email string, update *model.User) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	filter := bson.D{{Key: "email", Value: email}}

// 	update = bson.D{{"$set", bson.D{{"email", "newemail@example.com"}}}}

// 	_, err := m.DB.Collection("users").FindOneAndUpdate(ctx, filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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

func (m *MongoDB) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: email, Value: email}}

	r := m.DB.Collection("users").FindOne(ctx, filter)
	if r.Err() != mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	if r.Err() != nil {
		return nil, r.Err()
	}

	u := &model.User{}

	if err := r.Decode(u); err != nil {
		return nil, err
	}

	return u, nil
}

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
