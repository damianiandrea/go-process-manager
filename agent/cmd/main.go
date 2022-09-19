package main

import (
	"log"
	"os"

	"github.com/damianiandrea/go-process-manager/agent/internal/server"
)

func main() {
	addr := getenv("SERVER_ADDR", ":7000")
	nats := getenv("NATS_ADDR", "nats://127.0.0.1:4222")
	rate := getenv("HEART_RATE", "5s")

	s, err := server.New(
		server.WithAddr(addr),
		server.WithNats(nats),
		server.WithHeartRate(rate),
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
