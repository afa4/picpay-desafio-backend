package routine

import (
	"errors"
	"fmt"

	"github.com/afa4/picpay-desafio-backend/src/dao"
	"github.com/afa4/picpay-desafio-backend/src/entity"
)

type TransferRoutine struct {
	mongoDAO        *dao.MongoDAO
	transferChannel *chan entity.Transfer
}

func NewTransferRoutine(mongoDAO *dao.MongoDAO, transferChannel *chan entity.Transfer) *TransferRoutine {
	return &TransferRoutine{
		mongoDAO:        mongoDAO,
		transferChannel: transferChannel,
	}
}

func (tr *TransferRoutine) Start() {
	go tr.transferRoutine()
}

func (tr *TransferRoutine) transferRoutine() {
	fmt.Println("Transfer routine started")
	for transferReq := range *tr.transferChannel {
		fmt.Println("Transfer received")
		err := tr.executeTransfer(transferReq)
		if err != nil {
			fmt.Printf("Error executing transfer: %s\n", err.Error())
		}
	}
}

func (tr *TransferRoutine) executeTransfer(transferReq entity.Transfer) error {
	if transferReq.Payee == transferReq.Payer {
		return errors.New("Payer and payee are the same")
	}

	payerBalance, err := tr.mongoDAO.GetAccountBalance(transferReq.Payer)
	if err != nil {
		return err
	}

	if payerBalance.Balance < transferReq.Amount {
		// todo: could be a notification to payer
		return errors.New(fmt.Sprintf("Insufficient funds in payer account %d\n", transferReq.Payer))
	}

	debitTransaction := entity.Transaction{
		Type:             "debit",
		Amount:           transferReq.Amount,
		RelatedAccountID: transferReq.Payee,
	}

	creditTransaction := entity.Transaction{
		Type:             "credit",
		Amount:           transferReq.Amount,
		RelatedAccountID: transferReq.Payer,
	}

	//todo: handle rollback problems
	tr.mongoDAO.SaveTransaction(debitTransaction, transferReq.Payer)
	tr.mongoDAO.SaveTransaction(creditTransaction, transferReq.Payee)

	fmt.Printf("Transaction made: payerId=%d payeeId=%d\n", transferReq.Payer, transferReq.Payee)
	return nil
}
