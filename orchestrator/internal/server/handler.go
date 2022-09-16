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

type execProcessHandler struct {
	execProcessMsgProducer message.ExecProcessMsgProducer
}

func newExecProcessHandler(producer message.ExecProcessMsgProducer) *execProcessHandler {
	return &execProcessHandler{execProcessMsgProducer: producer}
}

func (h *execProcessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.execProcess(w, r)
	default:
		setAllowHeader(w, http.MethodOptions, http.MethodPost)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		writeJsonError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	}
}

func (h *execProcessHandler) execProcess(w http.ResponseWriter, r *http.Request) {
	request := &ExecProcessRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		writeJsonError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	execProcessMsg := &message.ExecProcess{ProcessName: request.ProcessName}
	if err := h.execProcessMsgProducer.Produce(r.Context(), execProcessMsg); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type ExecProcessRequest struct {
	ProcessName string `json:"process_name"`
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
