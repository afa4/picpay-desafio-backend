package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDAO struct {
	mongoClient *mongo.Client
}

func NewMongoDAO(uri string) *MongoDAO {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to mongo")
	}
	return &MongoDAO{mongoClient: client}
}

func (m *MongoDAO) GetTransactions(accountId int) ([]Transfer, error) {
	result, err := m.mongoClient.Database("picpay").Collection("transactions").Find(context.TODO(), bson.M{"payee": accountId})
	if err != nil {
		return nil, err
	}
	var transactions []Transfer
	for result.Next(context.TODO()) {
		var transaction Transfer
		err := result.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (m *MongoDAO) SaveTransaction(transaction Transfer) error {
	_, err := m.mongoClient.Database("picpay").Collection("transactions").InsertOne(context.TODO(), transaction)
	if err != nil {
		return err
	}
	return nil
}
