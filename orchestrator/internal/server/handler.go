package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
	"github.com/damianiandrea/go-process-manager/orchestrator/internal/storage"
)

var ErrInvalidRequestBody = errors.New("invalid request body")

type runProcessHandler struct {
	runProcessMsgProducer message.RunProcessMsgProducer
}

func newRunProcessHandler(producer message.RunProcessMsgProducer) *runProcessHandler {
	return &runProcessHandler{runProcessMsgProducer: producer}
}

func (h *runProcessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

type listRunningProcessesHandler struct {
	processStore storage.ProcessStore
}

func newListRunningProcessesHandler(processStore storage.ProcessStore) *listRunningProcessesHandler {
	return &listRunningProcessesHandler{processStore: processStore}
}

func (h *listRunningProcessesHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	runningProcesses, err := h.processStore.GetAll()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
		return
	}

	processes := make([]process, 0)
	for _, p := range runningProcesses {
		processes = append(processes, process{Pid: p.Pid, ProcessUuid: p.ProcessUuid, AgentId: p.AgentId,
			LastSeen: p.LastSeen})
	}
	response := listRunningProcessesResponse{Data: processes}
	writeJson(w, http.StatusOK, response)
}

type listRunningProcessesResponse struct {
	Data []process `json:"data"`
}

type process struct {
	Pid         int    `json:"pid"`
	ProcessUuid string `json:"process_uuid"`
	AgentId     string `json:"agent_id"`
	LastSeen    int64  `json:"last_seen"`
}

type processOutputHandler struct {
	upgrader    websocket.Upgrader
	msgConsumer message.ProcessOutputMsgConsumer
}

func newProcessOutputHandler(msgConsumer message.ProcessOutputMsgConsumer) *processOutputHandler {
	return &processOutputHandler{upgrader: websocket.Upgrader{}, msgConsumer: msgConsumer}
}

func (h *processOutputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	ctx := r.Context()
	outputCh := make(chan []byte, 64)
	defer close(outputCh)
	if err = h.msgConsumer.ChanConsume(ctx, mux.Vars(r)["processUuid"], outputCh); err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			case msg := <-outputCh:
				_ = conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
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
