package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

func writeJson(w http.ResponseWriter, code int, b interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(b)
}

func writeJsonError(w http.ResponseWriter, code int, err error) {
	response := errorResponse{Error: errorDetails{Code: code, Message: err.Error()}}
	writeJson(w, code, response)
}

func setAllowHeader(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", strings.Join(methods, ", "))
}

type errorResponse struct {
	Error errorDetails `json:"error"`
}

type errorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
