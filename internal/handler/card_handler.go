package handler

import (
	"bank-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

// CreateCard — POST /cards
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID  int    `json:"account_id"`
		CardHolder string `json:"card_holder"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	card, err := h.cardService.GenerateCard(req.AccountID, req.CardHolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// GetCards — GET /cards?account_id=...
func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	if accountIDStr == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)
		return
	}
	accountID, _ := strconv.Atoi(accountIDStr)

	cards, err := h.cardService.GetCardsByAccount(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cards)
}
