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

const (
	Url = "nats://nats-svc:4222"
	QueueGroup = "com-identity-queue"
	ClusterId = "devpie-client"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var clientID = fmt.Sprintf("com-identity-%d", rand.Int())

var repo, repoErr = NewRepository(Config{
	User:       os.Getenv("POSTGRES_USER"),
	Host:       os.Getenv("POSTGRES_HOST"),
	Name:       os.Getenv("POSTGRES_DB"),
	Password:   os.Getenv("POSTGRES_PASSWORD"),
	DisableTLS: true,
})

var infolog = log.New(os.Stdout, fmt.Sprintf("%s-identity-component: ", ClusterId), log.Lmicroseconds|log.Lshortfile)

func main() {
	if repoErr != nil {
		log.Fatal(repoErr)
	}
	defer repo.Close()

	infolog.Printf("main : Started")


	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Fatal(err)
	}
	AckWaitTimeout := stan.AckWait(dur)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	c, close := events.NewClient(ClusterId, clientID , Url)
	defer close()

	c.Listen(string(events.CommandsAddUser), QueueGroup, handleAddUserCommand, stan.DeliverAllAvailable(), stan.SetManualAckMode(),
		AckWaitTimeout, stan.DurableName(QueueGroup))

	select {
	case sig := <-shutdown:
		infolog.Println("main : Start shutdown", sig)
		close()
		infolog.Println("main : Closed NATS connection", sig)
		os.Exit(1)
	}
}
