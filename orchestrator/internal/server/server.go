package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message"
	"github.com/damianiandrea/go-process-manager/orchestrator/internal/message/nats"
	"github.com/damianiandrea/go-process-manager/orchestrator/internal/storage"
	"github.com/damianiandrea/go-process-manager/orchestrator/internal/storage/inmem"
)

var ErrNoMsgPlatform = errors.New("no message platform")

type server struct {
	ctx    context.Context
	stop   context.CancelFunc
	addr   string
	server *http.Server

	processStore storage.ProcessStore

	natsClient                      *nats.Client
	runProcessMsgProducer           message.RunProcessMsgProducer
	listRunningProcessesMsgConsumer message.ListRunningProcessesMsgConsumer
	processOutputMsgConsumer        message.ProcessOutputMsgConsumer
}

func New(options ...Option) (*server, error) {
	s := &server{}

	for _, opt := range options {
		opt(s)
	}

	s.processStore = inmem.NewProcessStore()
	encoder := &message.JsonEncoder{}
	decoder := &message.JsonDecoder{}

	if s.natsClient != nil {
		s.runProcessMsgProducer = nats.NewRunProcessMsgProducer(s.natsClient, encoder)
		s.listRunningProcessesMsgConsumer = nats.NewListRunningProcessesMsgConsumer(s.natsClient, decoder, s.processStore)
		s.processOutputMsgConsumer = nats.NewProcessOutputMsgConsumer(s.natsClient)
	} else {
		return nil, ErrNoMsgPlatform
	}

	r := mux.NewRouter()
	r.Use(recoverer)
	r.Path("/health").Methods(http.MethodGet).Handler(&healthHandler{})

	r.Path("/v1/processes").Methods(http.MethodPost).Handler(newRunProcessHandler(s.runProcessMsgProducer))
	r.Path("/v1/processes").Methods(http.MethodGet).Handler(newListRunningProcessesHandler(s.processStore))
	r.Path("/v1/processes/{processUuid}/logs").Methods(http.MethodGet).Handler(newProcessOutputHandler(s.processOutputMsgConsumer))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
	s.ctx = ctx
	s.stop = stop
	return s, nil
}

func (s *server) Run() error {
	defer s.cleanup()
	group, groupCtx := errgroup.WithContext(s.ctx)

	group.Go(func() error {
		log.Printf("listening on %v\n", s.server.Addr)
		return s.server.ListenAndServe()
	})

	group.Go(func() error {
		log.Println("ready to consume messages...")
		return s.listRunningProcessesMsgConsumer.Consume(groupCtx)
	})

	group.Go(func() error {
		<-groupCtx.Done()
		log.Println("gracefully shutting down...")
		return s.server.Shutdown(context.Background())
	})

	return group.Wait()
}

func (s *server) cleanup() {
	if s.natsClient != nil {
		s.natsClient.Close()
	}
	s.stop()
}

type Option func(*server)

func WithAddr(addr string) Option {
	return func(s *server) {
		s.addr = addr
	}
}

func WithNats(addr string) Option {
	return func(s *server) {
		client, err := nats.NewClient(addr)
		if err != nil {
			panic(err)
		}
		s.natsClient = client
	}
}
