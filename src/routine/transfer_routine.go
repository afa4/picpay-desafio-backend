package routine

import (
	"fmt"
	"time"

	"github.com/afa4/picpay-desafio-backend/src/dao"
	"github.com/afa4/picpay-desafio-backend/src/entity"
)

type TransferRoutine struct {
	mongoDAO        *dao.MongoDAO
	transferChannel *chan entity.Transfer
}

func NewTransferRoutine(mongoDAO dao.MongoDAO, transferChannel *chan entity.Transfer) *TransferRoutine {
	return &TransferRoutine{
		mongoDAO:        &mongoDAO,
		transferChannel: transferChannel,
	}
}

func (tr *TransferRoutine) Start() {
	go tr.transferRoutine()
}

func (tr *TransferRoutine) transferRoutine() {
	fmt.Println("Transfer routine started")
	for transferReq := range *tr.transferChannel {
		time.Sleep(1 * time.Second)
		fmt.Println("Transfer received")
		tr.executeTransfer(transferReq)
	}
}

func (tr *TransferRoutine) executeTransfer(transferReq entity.Transfer) {
	fmt.Println(transferReq)
	// todo:
	// get all payer transactions to check balance
	// sum transactions and get customer balance
	// register money transaction
	// notify receiver
}
