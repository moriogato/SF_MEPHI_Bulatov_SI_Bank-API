package handler

import (
	"bank-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// CreateAccount — POST /accounts
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	var req struct {
		Currency string `json:"currency"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	account, err := h.accountService.CreateAccount(userID, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetUserAccounts — GET /accounts
func (h *AccountHandler) GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	accounts, err := h.accountService.GetUserAccounts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

// Deposit — POST /accounts/{id}/deposit
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.accountService.Deposit(accountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deposit successful"})
}

// Withdraw — POST /accounts/{id}/withdraw
func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.accountService.Withdraw(accountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "withdrawal successful"})
}
