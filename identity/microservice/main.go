package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivorscott/devpie-client-events/go/events"

stan "github.com/nats-io/stan.go"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))


func main() {
	
    // ========================================
	// Required Variables
    // ========================================

	url := os.Getenv("NATS_URL")
	clu := os.Getenv("CLUSTER_ID")
	qg := fmt.Sprintf("%s-queue", os.Getenv("CLIENT_ID"))
	cid := fmt.Sprintf("%s-%d", os.Getenv("CLIENT_ID"), rand.Int())

	// ========================================
	// Default Logger Setup
    // ========================================

	infolog := log.New(os.Stdout, fmt.Sprintf("%s: ", cid), log.Lmicroseconds|log.Lshortfile)

	infolog.Printf("main : Started")
	
    // ========================================
	// Dedicated Database For Microservice
    // ========================================

	repo, err := NewRepository(Config{
		User:       os.Getenv("POSTGRES_USER"),
		Host:       os.Getenv("POSTGRES_HOST"),
		Name:       os.Getenv("POSTGRES_DB"),
		Password:   os.Getenv("POSTGRES_PASSWORD"),
		DisableTLS: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

    // ========================================
	// Create NATS Client
    // ========================================

	c, Close := events.NewClient(clu, cid, url)
	defer Close()

    // ========================================
	// Handle Messages
    // ========================================

	h := NewHandlers(repo, c)

	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Fatal(err)
	}

	streamName := fmt.Sprintf("%s.command",events.Identity)

	c.Listen(streamName, qg, h.handleAddUser, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(dur), stan.DurableName(qg))

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
