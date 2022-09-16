package server

import (
	"net/http"
)

type healthHandler struct {
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	writeJson(w, http.StatusOK, &healthResponse{Status: UP})
}

type healthResponse struct {
	Status health `json:"status"`
}

type health string

const (
	UP   health = "UP"
	DOWN        = "DOWN"
)
