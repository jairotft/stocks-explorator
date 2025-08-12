package handlers

import (
	"encoding/json"
	"net/http"
	"stock/backend/pkg/engine"
)

func GetBasicRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	recommendation, err := engine.GetDBRecommendations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recommendation, err = engine.GetOpenAIRecommendations(recommendation.Stocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recommendation)
}

func GetAdvancedRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
