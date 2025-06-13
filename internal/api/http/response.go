package http

import (
	"encoding/json"
	"net/http"
	"time"
)

type AccountResponse struct {
	AccountID int    `json:"account_id"`
	Balance   string `json:"balance"`
}

type TransactionResponse struct {
	TransactionID int64     `json:"transaction_id"`
	SourceID      int       `json:"source_account_id"`
	DestID        int       `json:"destination_account_id"`
	Amount        string    `json:"amount"`
	Timestamp     time.Time `json:"created_at"`
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
