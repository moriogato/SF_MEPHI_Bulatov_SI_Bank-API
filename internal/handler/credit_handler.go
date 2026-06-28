package handler

import (
	"bank-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	creditService *service.CreditService
}

func NewCreditHandler(creditService *service.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

// CreateCredit — POST /credits
func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	var req struct {
		AccountID  int     `json:"account_id"`
		Amount     float64 `json:"amount"`
		TermMonths int     `json:"term_months"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	credit, schedule, err := h.creditService.CreateCredit(userID, req.AccountID, req.Amount, req.TermMonths)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"credit":   credit,
		"schedule": schedule,
	})
}

// GetUserCredits — GET /credits
func (h *CreditHandler) GetUserCredits(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	credits, err := h.creditService.GetUserCredits(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credits)
}

// GetSchedule — GET /credits/{id}/schedule
func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creditID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid credit id", http.StatusBadRequest)
		return
	}

	schedule, err := h.creditService.GetSchedule(creditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(schedule)
}
