package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IMongoAdapter[T any] interface {
	FindAll(context context.Context, database string, collection string) (*[]T, error)
}

type MongoAdapter[T any] struct {
	client *mongo.Client
}

func NewMongoMongoAdapter[T any](client *mongo.Client) *MongoAdapter[T] {
	return &MongoAdapter[T]{
		client: client,
	}
}

func (m *MongoAdapter[T]) FindAll(context context.Context, database string, collection string) (*[]T, error) {
	result, err := m.client.Database(database).Collection(collection).Find(context, bson.D{})
	if err != nil {
		return nil, err
	}
	var list []T
	for result.Next(context) {
		var item T
		err := result.Decode(&item)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return &list, nil
}
