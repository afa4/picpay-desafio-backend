package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/afa4/picpay-desafio-backend/src/dao"
)

type AccountBalanceController struct {
	mongoDAO *dao.MongoDAO
}

func NewAccountBalanceController(mongoDAO *dao.MongoDAO) *AccountBalanceController {
	return &AccountBalanceController{
		mongoDAO: mongoDAO,
	}
}

func (ac *AccountBalanceController) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ac.GetAccountBalance(w, r)
	}
}

func (ac *AccountBalanceController) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.URL.Query().Get("account_id"))
	if err != nil {
		http.Error(w, "account_id (integer) is required", http.StatusBadRequest)
		return
	}

	balance, err := ac.mongoDAO.GetAccountBalance(accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, "error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
