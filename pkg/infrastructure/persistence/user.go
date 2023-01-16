package persistence

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"golang-jwt-example/pkg/domain/repository"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collection = "users"

type UserRepository struct {
	database *mongo.Collection
}

var _ repository.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		database: db.Collection(collection),
	}
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*entity.User, error) {
	users := make([]*entity.User, 0)
	srt := bson.D{
		primitive.E{Key: "_id", Value: 1},
	}
	opt := options.Find().SetSort(srt)
	flt := bson.D{
		primitive.E{},
	}
	cur, err := r.database.Find(ctx, flt, opt)
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &users); err != nil {
		log.Printf("error %+v", err)
		return nil, err
	}

	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (r *UserRepository) CreateUser(ctx context.Context, inputData input.User) (interface{}, error) {
	// opts := options.InsertOneOptions{}
	data := input.User{
		Name:     inputData.Name,
		Password: inputData.Password,
	}
	id, err := r.database.InsertOne(ctx, data, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return id, nil
}
