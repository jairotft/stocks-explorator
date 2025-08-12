package handlers

import (
	"encoding/json"
	"net/http"

	"stock/backend/pkg/engine"
	"stock/backend/pkg/exceptions"
)

func GetStocksHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	brokerage := r.URL.Query().Get("brokerage")
	action := r.URL.Query().Get("action")
	ratingFrom := r.URL.Query().Get("rating_from")
	ratingTo := r.URL.Query().Get("rating_to")
	page := r.URL.Query().Get("page")
	ticker := r.URL.Query().Get("ticker")
	mapParams := map[string]string{
		"ticker":      ticker,
		"brokerage":   brokerage,
		"action":      action,
		"rating_from": ratingFrom,
		"rating_to":   ratingTo,
		"page":        page,
	}
	stocksResponse, err := engine.GetStocks(mapParams)
	if err != nil {
		exceptions.Throw(w, exceptions.AppException{Detail: "Error generando respuesta"}, http.StatusInternalServerError, err)
		return
	}

	payloadResponse := map[string]interface{}{
		"message": "Success",
		"data":    stocksResponse,
	}

	response, err := json.Marshal(payloadResponse)
	if err != nil {
		exceptions.Throw(w, exceptions.AppException{Detail: "Error generando respuesta"}, http.StatusInternalServerError, err)
		return
	}

	w.Write(response)
}
