package main

import (
	"fmt"
	"github.com/ivorscott/devpie-client-events/go/events"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	stan "github.com/nats-io/stan.go"
)

var SeededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const (
	Url = "nats://nats-svc:4222"
	QueueGroup = "com-identity-queue"
	ClusterId = "devpie-client"

)

func handleNewUser(m *stan.Msg) {
	fmt.Printf("Received a message: %s\n", string(m.Data))

	// TODO: save user in dedicated database

	m.Ack()
}

func main() {
	infolog := log.New(os.Stdout, "devpie-client-identity-component: ", log.Lmicroseconds|log.Lshortfile)

	infolog.Printf("main : Started")

	var ClientID = fmt.Sprintf("com-identity-%d", rand.Int())

	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Fatal(err)
	}
	AckWaitTimeout := stan.AckWait(dur)

	// Listen for an interrupt or terminate signal.
	// Signal package requires a buffered channel.

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	c, close := events.NewClient(ClusterId, ClientID , Url)
	defer close()

	c.Listen(string(events.CommandsAddUser), QueueGroup, handleNewUser, stan.DeliverAllAvailable(), stan.SetManualAckMode(),
		AckWaitTimeout, stan.DurableName(QueueGroup))

	select {
	// Graceful shutdown
	case sig := <-shutdown:
		infolog.Println("main : Start shutdown", sig)
		close()
		infolog.Println("main : Closed NATS connection", sig)
		os.Exit(1)
	}
}
