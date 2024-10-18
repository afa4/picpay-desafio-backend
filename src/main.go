package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/afa4/picpay-desafio-backend/src/controller"
	dao "github.com/afa4/picpay-desafio-backend/src/dao"
	entity "github.com/afa4/picpay-desafio-backend/src/entity"
)

func main() {
	mongoDAO := dao.NewMongoDAO("mongodb://root:example@localhost:27017/")
	transferChannel := make(chan entity.Transfer)
	transferController := controller.NewTransferController(mongoDAO, &transferChannel)
	rootController := controller.NewRootController()
	go transferRoutine(&transferChannel, mongoDAO)
	http.HandleFunc("/transfer", transferController.HandlerFunc())
	http.HandleFunc("/", rootController.HandlerFunc())
	http.ListenAndServe(":8080", nil)
}

func transferRoutine(transferChannel *chan entity.Transfer, mongoDAO *dao.MongoDAO) {
	fmt.Println("Transfer routine started")
	for transferReq := range *transferChannel {
		time.Sleep(1 * time.Second)
		fmt.Println("Transfer received")
		executeTransfer(transferReq, mongoDAO)
	}
}

func executeTransfer(transferReq entity.Transfer, mongoDAO *dao.MongoDAO) {
	fmt.Println(transferReq)
	// get all payer transactions to check balance
	// sum transactions and get customer balance
	// register money transaction
	// notify receiver
}

// todo: use channel to process Transfer atomically
