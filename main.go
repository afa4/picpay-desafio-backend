package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Transaction struct {
	Payer  int32   `json:"payer"`
	Payee  int32   `json:"payee"`
	Amount float64 `json:"value"`
}

func main() {
	transactionChannel := make(chan Transaction)
	go transactionRoutine(&transactionChannel)
	http.HandleFunc("/", getRootHandler())
	http.HandleFunc("/transfer", postTransactionHandler(&transactionChannel))
	http.ListenAndServe(":8080", nil)
}

func getRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Hello World!"))
	}
}

func postTransactionHandler(transactionChannel *chan Transaction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		transaction := Transaction{}
		json.Unmarshal(body, &transaction)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		*transactionChannel <- transaction
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Hello World!" + " " + fmt.Sprintf("%d", transaction.Payer) + " " + fmt.Sprintf("%d", transaction.Payee) + " " + fmt.Sprintf("%f", transaction.Amount)))
	}
}

func transactionRoutine(transactionChannel *chan Transaction) {
	fmt.Println("Transaction routine started")
	for {
		transaction := <-*transactionChannel
		time.Sleep(1 * time.Second)
		fmt.Println("Transaction received")
		fmt.Println(transaction)
	}
}

// todo: use channel to process transaction atomically
