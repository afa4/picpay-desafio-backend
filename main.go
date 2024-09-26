package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func main() {
	mongoDAO := NewMongoDAO("mongodb://root:example@mongo:27017/")
	transferChannel := make(chan Transfer)
	go transferRoutine(&transferChannel, mongoDAO)
	http.HandleFunc("/transfer", handleTransfer(&transferChannel, mongoDAO))
	http.HandleFunc("/", handleRoot())
	http.ListenAndServe(":8080", nil)
}

func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	}
}

func handleTransfer(transferChannel *chan Transfer, mongoDAO *MongoDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			postTransferHandler(w, r, transferChannel)
			return
		}
		if r.Method == "GET" {
			getTransferHandler(w, r, mongoDAO)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func postTransferHandler(w http.ResponseWriter, r *http.Request, transferChannel *chan Transfer) {
	body, err := io.ReadAll(r.Body)
	transferReq := Transfer{}
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &transferReq)
	if err != nil {
		http.Error(w, "Error parsing json", http.StatusInternalServerError)
		return
	}
	*transferChannel <- transferReq
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Hello World!" + " " + fmt.Sprintf("%d", transferReq.Payer) + " " + fmt.Sprintf("%d", transferReq.Payee) + " " + fmt.Sprintf("%f", transferReq.Amount)))
}

func getTransferHandler(w http.ResponseWriter, r *http.Request, mongoDAO *MongoDAO) {
	accountIdStr := r.URL.Query().Get("account_id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	transfers, err := mongoDAO.GetTransactions(accountId)
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

func transferRoutine(transferChannel *chan Transfer, mongoDAO *MongoDAO) {
	fmt.Println("Transfer routine started")
	for transferReq := range *transferChannel {
		time.Sleep(1 * time.Second)
		fmt.Println("Transfer received")
		executeTransfer(transferReq, mongoDAO)
	}
}

func executeTransfer(transferReq Transfer, mongoDAO *MongoDAO) {
	fmt.Println(transferReq)
	// get all payer transactions to check balance
	// sum transactions and get customer balance
	// register money transaction
	// notify receiver
}

// todo: use channel to process Transfer atomically
