package dao

import (
	"context"
	"fmt"

	entity "github.com/afa4/picpay-desafio-backend/src/entity"
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

func (m *MongoDAO) GetAccountBalance(accountId int) (*entity.Balance, error) {
	transactions, err := m.GetTransactions(accountId)

	if err != nil {
		return nil, err
	}

	balance := 0.0
	for _, transaction := range transactions {
		if transaction.Type == "credit" {
			balance += transaction.Amount
		} else {
			balance -= transaction.Amount
		}
	}

	return &entity.Balance{Balance: balance}, nil
}

func (m *MongoDAO) GetTransactions(accountId int) ([]entity.Transaction, error) {
	srtAccountId := fmt.Sprintf("acc_%d", accountId)
	result, err := m.mongoClient.Database("transactions").Collection(srtAccountId).Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	var transactions []entity.Transaction
	for result.Next(context.TODO()) {
		var transaction entity.Transaction
		err := result.Decode(&transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (m *MongoDAO) SaveTransaction(transaction entity.Transaction, accountId int) error {
	srtAccountId := fmt.Sprintf("acc_%d", accountId)
	_, err := m.mongoClient.Database("transactions").Collection(srtAccountId).InsertOne(context.TODO(), transaction)
	if err != nil {
		return err
	}
	return nil
}
