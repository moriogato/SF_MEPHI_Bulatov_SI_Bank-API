package handler

import (
	"bank-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(transferService *service.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

// Transfer — POST /transfer
func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID int     `json:"from_account_id"`
		ToAccountID   int     `json:"to_account_id"`
		Amount        float64 `json:"amount"`
		Description   string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	tx, err := h.transferService.Transfer(req.FromAccountID, req.ToAccountID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

// GetTransactions — GET /transactions?account_id=...
func (h *TransferHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)
		return
	}
	accountID, _ := strconv.Atoi(accountIDStr)

	transactions, err := h.transferService.GetAccountTransactions(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
