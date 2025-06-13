package http

import (
	nativehttp "net/http"

	chi "github.com/go-chi/chi/v5"
)

func NewRouter(h *Handler) nativehttp.Handler {
	r := chi.NewRouter()
	r.Post("/accounts", h.CreateAccount)
	r.Get("/accounts/{account_id}", h.GetAccount)
	r.Post("/transactions", h.SubmitTransaction)
	r.Post("/deposit", h.Deposit)
	r.Post("/withdraw", h.Withdraw)
	return r
}
