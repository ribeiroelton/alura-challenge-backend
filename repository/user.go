package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type UserRepository struct {
	db *mongo.Database
}

const (
	userCollection = "users"
)

func NewUserRepository(c *config.Config) (*UserRepository, error) {
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

	return &UserRepository{
		db: db,
	}, nil
}

//SaveUser saves a new user to the database if the specified email does not already exists.
func (ur *UserRepository) SaveUser(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := ur.db.Collection(userCollection).InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

//UpdateUserByEmail updates an user using its email as key. Default timeout of 3 seconds
//TODO implements DeleteUserByEmail
func (r *UserRepository) UpdateUserByEmail(update *model.User) error {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return nil
}

//DeleteUserByEmail deletes an user by using its email. Default timeout of 3 seconds
func (ur *UserRepository) DeleteUserByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: email, Value: email}}

	r, err := ur.db.Collection(userCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if r.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

//GetUserByEmail returns an user by its email. Default timeout of 3 seconds
func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: email, Value: email}}

	r := ur.db.Collection(userCollection).FindOne(ctx, filter)
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
func (ur *UserRepository) HasUserByEmail(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}

	r := ur.db.Collection(userCollection).FindOne(ctx, filter)
	u := model.User{}
	r.Decode(&u)
	if r.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if r.Err() != nil {
		return false, r.Err()
	}

	return true, nil
}

//ListUsers returns a slice of users. Empty if no users found. Default timeout of 30 seconds.
func (ur *UserRepository) ListUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	users := []model.User{}

	cur, err := ur.db.Collection(userCollection).Find(ctx, bson.D{{}})
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
