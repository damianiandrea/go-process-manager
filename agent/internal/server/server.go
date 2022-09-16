package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
	"github.com/damianiandrea/go-process-manager/agent/internal/message/nats"
	"github.com/damianiandrea/go-process-manager/agent/internal/process/local"
)

var ErrNoMsgConsumer = errors.New("no message consumer")

type server struct {
	ctx    context.Context
	stop   context.CancelFunc
	addr   string
	server *http.Server

	natsClient  *nats.Client
	msgConsumer message.ExecProcessMsgConsumer
}

func New(options ...Option) (*server, error) {
	s := &server{}

	for _, opt := range options {
		opt(s)
	}

	if s.natsClient != nil {
		s.msgConsumer = nats.NewExecProcessMsgConsumer(s.natsClient, &local.ProcessExecutor{})
	} else {
		return nil, ErrNoMsgConsumer
	}

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
		return s.msgConsumer.Consume(groupCtx)
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
