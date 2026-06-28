package handler

import (
	"bank-api/internal/service"
	"encoding/json"
	"net/http"
)

type CBRHandler struct {
	cbrService *service.CBRService
}

func NewCBRHandler(cbrService *service.CBRService) *CBRHandler {
	return &CBRHandler{cbrService: cbrService}
}

// GetKeyRate — GET /cbr/rate
func (h *CBRHandler) GetKeyRate(w http.ResponseWriter, r *http.Request) {
	rate, err := h.cbrService.GetKeyRate()
	if err != nil {
		http.Error(w, "Failed to get key rate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"key_rate": rate,
	})
}
