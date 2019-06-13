package controllers

import (
	"encoding/json"
	"net/http"
)

type BaseController struct{}

type JsonResponse struct {
	Error   int         `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}

// Writes the response as a standard JSON response with StatusOK
func (base *BaseController) sendOk(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&JsonResponse{Error: 0, Data: m, Message: ""}); err != nil {
		base.sendError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func (base *BaseController) sendError(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&JsonResponse{Error: errorCode, Data: "", Message: errorMsg})
}
