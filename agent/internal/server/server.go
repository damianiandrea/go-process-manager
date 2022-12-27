package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/message/nats"
	"github.com/damianiandrea/go-process-manager/agent/internal/process"
	"github.com/damianiandrea/go-process-manager/agent/internal/process/local"
	"github.com/damianiandrea/go-process-manager/agent/internal/scheduler"
)

var ErrNoMsgPlatform = errors.New("no message platform")

type server struct {
	ctx    context.Context
	stop   context.CancelFunc
	addr   string
	server *http.Server

	agentId        string
	processManager process.Manager

	natsClient                      *nats.Client
	runProcessMsgConsumer           message.RunProcessMsgConsumer
	listRunningProcessesMsgProducer message.ListRunningProcessesMsgProducer
	processOutputMsgProducer        message.ProcessOutputMsgProducer

	heartRate          time.Duration
	heartbeatScheduler *scheduler.Heartbeat
}

func New(options ...Option) (*server, error) {
	s := &server{}
	s.agentId = uuid.NewString()

	for _, opt := range options {
		opt(s)
	}

	encoder := &message.JsonEncoder{}
	decoder := &message.JsonDecoder{}

	if s.natsClient != nil {
		s.processOutputMsgProducer = nats.NewProcessOutputMsgProducer(s.agentId, s.natsClient)
		s.processManager = local.NewProcessManager(s.processOutputMsgProducer)
		s.runProcessMsgConsumer = nats.NewRunProcessMsgConsumer(s.natsClient, decoder, s.processManager)
		s.listRunningProcessesMsgProducer = nats.NewListRunningProcessesMsgProducer(s.agentId, s.natsClient, encoder)
	} else {
		return nil, ErrNoMsgPlatform
	}

	s.heartbeatScheduler = scheduler.NewHeartbeat(s.heartRate, s.processManager, s.listRunningProcessesMsgProducer)

	mux := http.NewServeMux()
	mux.Handle("/health", recoverer(&healthHandler{}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
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
		return s.runProcessMsgConsumer.Consume(groupCtx)
	})

	group.Go(func() error {
		log.Println("heart is beating...")
		return s.heartbeatScheduler.Beat(groupCtx)
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

func WithHeartRate(rate string) Option {
	return func(s *server) {
		parsed, err := time.ParseDuration(rate)
		if err != nil {
			panic(err)
		}
		s.heartRate = parsed
	}
}
