package subs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"wb-l0/internal/cache"
	"wb-l0/internal/model"
	"wb-l0/internal/repo"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type Subs struct {
	repo repo.Repo
	cache *cache.Cache
}

func (subs *Subs) printMsg(m *stan.Msg) {
	var order model.Order

	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		log.Printf("Error adding to db: %s", err)
	}
	subs.repo.AddToDB(order)
	subs.cache.AddToCache(order)
}

func (subs *Subs) InitSubs() {

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
		subs.printMsg(msg)
	}
	sub, err := sc.Subscribe(subj, mcb);
	//sub, err := sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}
	// TODO: поправить
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