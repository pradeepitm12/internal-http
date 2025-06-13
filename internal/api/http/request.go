package http

type CreateAccountRequest struct {
	AccountID      int    `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type TransferRequest struct {
	SourceAccountID      int    `json:"source_account_id"`
	DestinationAccountID int    `json:"destination_account_id"`
	Amount               string `json:"amount"`
}

type AccountOperationRequest struct {
	AccountID int    `json:"account_id"`
	Amount    string `json:"amount"`
}
