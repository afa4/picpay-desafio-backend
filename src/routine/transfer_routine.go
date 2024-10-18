package routine

import (
	"fmt"

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
		fmt.Println("Transfer received")
		tr.executeTransfer(transferReq)
	}
}

func (tr *TransferRoutine) executeTransfer(transferReq entity.Transfer) {
	if transferReq.Payee == transferReq.Payer {
		fmt.Println("Payer and payee are the same")
		return
	}

	payerBalance := tr.getBalance(transferReq.Payer)
	fmt.Printf("account %d balance $%f\n", transferReq.Payer, payerBalance)

	if payerBalance < transferReq.Amount {
		// todo: could be a notification to payer
		fmt.Printf("Insufficient funds in payer account %d\n", transferReq.Payer)
		return
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
}

func (tr *TransferRoutine) getBalance(accountId int) float64 {
	transactions, err := tr.mongoDAO.GetTransactions(accountId)

	if err != nil {
		fmt.Println("Error getting transactions")
		return 0
	}

	balance := 0.0
	for _, transaction := range transactions {
		if transaction.Type == "credit" {
			balance += transaction.Amount
		} else {
			balance -= transaction.Amount
		}
	}

	return balance
}
