package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var usageStr = `
Usage: stan-sub [options] <subject>

Options:
	-s,  --server   <url>            NATS Streaming server URL(s)
	-c,  --cluster  <cluster name>   NATS Streaming cluster name
	-id, --clientid <client ID>      NATS Streaming client ID
	-cr, --creds    <credentials>    NATS 2.0 Credentials

Subscription Options:
	--qgroup <name>                  Queue group
	--all                            Deliver all available messages
	--last                           Deliver starting with last published message
	--since  <time_ago>              Deliver messages in last interval (e.g. 1s, 1hr)
	--seq    <seqno>                 Start at seqno
	--new_only                       Only deliver new messages
	--durable <name>                 Durable subscriber name
	--unsub                          Unsubscribe the durable on exit
`

// NOTE: Use tls scheme for TLS, e.g. stan-sub -s tls://demo.nats.io:4443 foo
func usage() {
	log.Fatalf(usageStr)
}

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

func main() {

	// Connect to NATS
	nc, err := nats.Connect("localhost:4222")	// стандартный порт nats
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Connect to NATS-STREAMING
	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, "localhost:4222")
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", "localhost:4222", "test-cluster", "test-client")

	subj, i := "test-subj", 0	// тема сообщения
	mcb := func(msg *stan.Msg) {
		i++
		printMsg(msg, i)
	}
	sub, err := sc.Subscribe(subj, mcb);
	//sub, err := sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", subj)

	// 
	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)	// Канал - механизм общения горутинами
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)	// os.Interrupt = ctrl + c
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone	// булевый канал ждет когда из него что-то можно достать
}