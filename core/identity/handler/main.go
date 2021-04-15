package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devpies/devpie-client-events/go/events"
	stan "github.com/nats-io/stan.go"
)

func main() {
	url := os.Getenv("NATS_URL")
	clu := os.Getenv("CLUSTER_ID")
	cid := fmt.Sprintf("%s-%d", os.Getenv("CLIENT_ID"), rand.Int())

	infolog := log.New(os.Stdout, fmt.Sprintf("%s: ", cid), log.Lmicroseconds|log.Lshortfile)

	infolog.Printf("main : Started")

	// ========================================
	// Create NATS Client
	// ========================================

	c, Close := events.NewClient(clu, cid, url)
	defer Close()

	// ========================================
	// Handle Messages
	// ========================================

	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Fatal(err)
	}

	h := NewHandlers(c, infolog)

	streamName := fmt.Sprintf("%s.command", events.Identity)

	c.Listen(streamName, h.handleAddUser, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(dur))

	// ========================================
	// Graceful Shutdown
	// ========================================

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-shutdown:
		infolog.Println("main : Start shutdown", sig)
		Close()
		infolog.Println("main : Closed NATS connection", sig)
		os.Exit(1)
	}
}
