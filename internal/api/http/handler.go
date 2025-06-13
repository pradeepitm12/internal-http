package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pradeepitm12/compaaa/internal-http/internal/domain/model"
	"github.com/pradeepitm12/compaaa/internal-http/internal/transfer"
	"github.com/pradeepitm12/compaaa/internal-http/util"
)

type Handler struct {
	transferService *transfer.Service
	accountRepo     transfer.AccountRepository
}

// NewHandler wires the services and returns a handler instance.
func NewHandler(ts *transfer.Service, ar transfer.AccountRepository) *Handler {
	return &Handler{
		transferService: ts,
		accountRepo:     ar,
	}
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	val, err := util.StringToDecimal(req.InitialBalance)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid initial balance"), http.StatusBadRequest)
	}
	acc := model.NewAccount(req.AccountID, val)

	h.transferService.TxManager.Do(r.Context(), func(ctx context.Context, tx *sql.Tx) error {
		if err := h.accountRepo.Create(r.Context(), tx, acc); err != nil {
			fmt.Printf("error updating account %v\n", err)
			http.Error(w, "failed to create account", http.StatusInternalServerError)
		}
		return nil
	})

	fmt.Println("account created")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("account_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}
	resp := AccountResponse{}
	h.transferService.TxManager.Do(r.Context(), func(ctx context.Context, tx *sql.Tx) error {
		acc, err := h.accountRepo.GetByID(r.Context(), tx, id)
		if err != nil {
			http.Error(w, "account not found", http.StatusNotFound)
		}
		resp.AccountID = acc.ID()
		resp.Balance = acc.Balance().String()
		return nil
	})

	writeJSON(w, resp, http.StatusOK)
}

func (h *Handler) SubmitTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	amount, err := util.StringToDecimal(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	tx, err := h.transferService.Transfer(r.Context(), req.SourceAccountID, req.DestinationAccountID, amount)
	if err != nil {
		http.Error(w, "transfer failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	resp := TransactionResponse{
		TransactionID: tx.ID(),
		SourceID:      tx.SourceAccountID(),
		DestID:        tx.DestAccountID(),
		Amount:        tx.Amount().String(),
		Timestamp:     tx.CreatedAt(),
	}
	writeJSON(w, resp, http.StatusCreated)
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req AccountOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	amount, err := util.StringToDecimal(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	acc, err := h.transferService.Deposit(r.Context(), req.AccountID, amount)
	if err != nil {
		http.Error(w, "deposit failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	writeJSON(w, AccountResponse{
		AccountID: acc.ID(),
		Balance:   acc.Balance().String(),
	}, http.StatusOK)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var req AccountOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	amount, err := util.StringToDecimal(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	acc, err := h.transferService.Withdraw(r.Context(), req.AccountID, amount)
	if err != nil {
		http.Error(w, "withdraw failed: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	writeJSON(w, AccountResponse{
		AccountID: acc.ID(),
		Balance:   acc.Balance().String(),
	}, http.StatusOK)
}
