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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		transaction := Transaction{}
		json.Unmarshal(body, &transaction)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Hello World!" + " " + fmt.Sprintf("%d", transaction.Payer) + " " + fmt.Sprintf("%d", transaction.Payee) + " " + fmt.Sprintf("%f", transaction.Amount)))
	})
	http.ListenAndServe(":8080", nil)
}

// todo: use channel to process transaction atomically
