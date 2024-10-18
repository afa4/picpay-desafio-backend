package main

import (
	"net/http"

	"github.com/afa4/picpay-desafio-backend/src/controller"
	dao "github.com/afa4/picpay-desafio-backend/src/dao"
	entity "github.com/afa4/picpay-desafio-backend/src/entity"
	"github.com/afa4/picpay-desafio-backend/src/routine"
)

func main() {
	// dependencies initialization
	mongoDAO := dao.NewMongoDAO("mongodb://root:example@localhost:27017/")
	transferChannel := make(chan entity.Transfer)

	// concurrent routine initialization
	transferRoutine := routine.NewTransferRoutine(*mongoDAO, &transferChannel)
	transferRoutine.Start()

	// controller initialization
	transferController := controller.NewTransferController(mongoDAO, &transferChannel)
	rootController := controller.NewRootController()

	// http server initialization
	http.HandleFunc("/transfer", transferController.HandlerFunc())
	http.HandleFunc("/", rootController.HandlerFunc())
	http.ListenAndServe(":8080", nil)
}
