package main

import (
	"log"
	"os"

	"github.com/damianiandrea/go-process-manager/orchestrator/internal/server"
)

func main() {
	addr := getenv("SERVER_ADDR", ":6000")
	nats := getenv("NATS_ADDR", "nats://127.0.0.1:4222")

	s, err := server.New(
		server.WithAddr(addr),
		server.WithNats(nats),
	)
	if err != nil {
		log.Fatalf("could not create server: %v\n", err)
	}

	log.Fatalf("exiting: %v\n", s.Run())
}

func getenv(key, def string) string {
	if env, found := os.LookupEnv(key); found {
		return env
	}
	return def
}
