package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/afa4/picpay-desafio-backend/src/dao"
	"github.com/afa4/picpay-desafio-backend/src/entity"
)

type TransferController struct {
	mongoDAO        dao.IMongoDAO
	transferChannel *chan entity.Transfer
	mongoAdapter    dao.IMongoAdapter[entity.Transaction]
}

func NewTransferController(mongoDAO *dao.MongoDAO, transferChannel *chan entity.Transfer, mongoAdapter dao.IMongoAdapter[entity.Transaction]) *TransferController {
	return &TransferController{
		mongoDAO:        mongoDAO,
		transferChannel: transferChannel,
		mongoAdapter:    mongoAdapter,
	}
}

func (tc *TransferController) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			tc.postTransferHandler(w, r)
		case "GET":
			tc.getTransferHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (tc *TransferController) postTransferHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	transferReq := entity.Transfer{}
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &transferReq)
	if err != nil {
		http.Error(w, "Error parsing json", http.StatusInternalServerError)
		return
	}
	*tc.transferChannel <- transferReq
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Hello World!" + " " + fmt.Sprintf("%d", transferReq.Payer) + " " + fmt.Sprintf("%d", transferReq.Payee) + " " + fmt.Sprintf("%f", transferReq.Amount)))
}

func (tc *TransferController) getTransferHandler(w http.ResponseWriter, r *http.Request) {
	accountIdStr := r.URL.Query().Get("account_id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	transfers, err := tc.mongoAdapter.FindAll(context.TODO(), "transactions", fmt.Sprintf("acc_%d", accountId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(transfers)
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
