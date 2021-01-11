package events

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
	"time"
)

type Client struct {
	Conn stan.Conn
}

func NewClient(clusterId, clientId, url string) (*Client, func () error) {
	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(url),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection failed for some reason: %v", reason)
		}))
	if err != nil {
		fmt.Printf("error: %v",err)
	}

	return &Client{
		Conn: sc,
	}, sc.Close
}

func (c *Client) Listen(subj, quegrp string, parseMsg func(msg *stan.Msg), opts ...stan.SubscriptionOption) {
	_, err := c.Conn.QueueSubscribe(subj, quegrp, parseMsg, opts...)
	if err != nil {
		c.Conn.Close()
		log.Fatal(err)
	}
}

func (c *Client) Publish(sub string, msg []byte) {
	ch := make(chan bool)
	var glock sync.Mutex
	var guid string

	acb := func(lguid string, err error) {
		glock.Lock()
		log.Printf("Received ACK for guid %s\n", lguid)
		defer glock.Unlock()
		if err != nil {
			log.Fatalf("Error in server ack for guid %s: %v\n", lguid, err)
		}
		if lguid != guid {
			log.Fatalf("Expected a matching guid in ack callback, got %s vs %s\n", lguid, guid)
		}
		ch <- true
	}

	glock.Lock()
	guid, err := c.Conn.PublishAsync(sub, msg, acb)
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
	glock.Unlock()
	if guid == "" {
		log.Fatal("Expected non-empty guid to be returned.")
	}
	log.Printf("Published [%s] : '%s' [guid: %s]\n", sub, msg, guid)

	select {
	case <-ch:
		break
	case <-time.After(5 * time.Second):
		log.Fatal("timeout")
	}
}
