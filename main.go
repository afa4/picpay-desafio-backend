package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TransferRequest struct {
	Payer  int32   `json:"payer"`
	Payee  int32   `json:"payee"`
	Amount float64 `json:"value"`
}

func main() {
	transferChannel := make(chan TransferRequest)
	go transferRoutine(&transferChannel)
	http.HandleFunc("/", getRootHandler())
	http.HandleFunc("/transfer", postTransferHandler(&transferChannel))
	http.ListenAndServe(":8080", nil)
}

func getRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	}
}

func postTransferHandler(transferChannel *chan TransferRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		transferReq := TransferRequest{}
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
}

func transferRoutine(transferChannel *chan TransferRequest) {
	fmt.Println("Transfer routine started")
	for transferReq := range *transferChannel {
		time.Sleep(1 * time.Second)
		fmt.Println("Transfer received")
		executeTransfer(transferReq)
	}
}

func executeTransfer(transferReq TransferRequest) {
	fmt.Println(transferReq)
}

// todo: use channel to process Transfer atomically
