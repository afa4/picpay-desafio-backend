package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/afa4/picpay-desafio-backend/src/dao"
	"github.com/afa4/picpay-desafio-backend/src/entity"
)

type TransferController struct {
	mongoDAO        *dao.MongoDAO
	transferChannel *chan entity.Transfer
}

func NewTransferController(mongoDAO *dao.MongoDAO, transferChannel *chan entity.Transfer) *TransferController {
	return &TransferController{
		mongoDAO:        mongoDAO,
		transferChannel: transferChannel,
	}
}

func (c *TransferController) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			c.postTransferHandler(w, r)
		case "GET":
			c.getTransferHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (c *TransferController) postTransferHandler(w http.ResponseWriter, r *http.Request) {
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
	*c.transferChannel <- transferReq
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Hello World!" + " " + fmt.Sprintf("%d", transferReq.Payer) + " " + fmt.Sprintf("%d", transferReq.Payee) + " " + fmt.Sprintf("%f", transferReq.Amount)))
}

func (c *TransferController) getTransferHandler(w http.ResponseWriter, r *http.Request) {
	accountIdStr := r.URL.Query().Get("account_id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	transfers, err := c.mongoDAO.GetTransactions(accountId)
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
