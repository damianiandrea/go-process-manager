package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
)

var (
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrInvalidRequestBody = errors.New("invalid request body")
)

type runProcessHandler struct {
	runProcessMsgProducer message.RunProcessMsgProducer
}

func newRunProcessHandler(producer message.RunProcessMsgProducer) *runProcessHandler {
	return &runProcessHandler{runProcessMsgProducer: producer}
}

func (h *runProcessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.runProcess(w, r)
	default:
		setAllowHeader(w, http.MethodOptions, http.MethodPost)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		writeJsonError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
}

func (h *runProcessHandler) runProcess(w http.ResponseWriter, r *http.Request) {
	request := &runProcessRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		writeJsonError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	runProcessMsg := &message.RunProcess{ProcessName: request.ProcessName, Args: request.Args}
	if err := h.runProcessMsgProducer.Produce(r.Context(), runProcessMsg); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type runProcessRequest struct {
	ProcessName string   `json:"process_name"`
	Args        []string `json:"args,omitempty"`
}

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
