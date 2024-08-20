package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Transaction struct {
	Payer  int32   `json:"payer"`
	Payee  int32   `json:"payee"`
	Amount float64 `json:"value"`
}

func main() {
	http.HandleFunc("/", getRootHandler())
	http.HandleFunc("/transfer", postTransactionHandler())
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

func postTransactionHandler() http.HandlerFunc {
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

		w.Write([]byte("Hello World!" + " " + fmt.Sprintf("%d", transaction.Payer) + " " + fmt.Sprintf("%d", transaction.Payee) + " " + fmt.Sprintf("%f", transaction.Amount)))
	}
}

// todo: use channel to process transaction atomically
