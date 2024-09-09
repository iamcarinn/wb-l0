package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"wb-l0/internal/cache"
	"wb-l0/internal/model"
	"wb-l0/internal/repo"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "pub"
	URL       = stan.DefaultNatsURL
	subj      = "test-subj"
)

func main() {
	// Создаем новый репозиторий
	r := repo.New()
	c := cache.New()

	// Чтение JSON-данных заказа из терминала
	fmt.Println("Enter JSON order string:")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading data from terminal: %v", err)
	}

	// Декодируем JSON в структуру Order
	var order model.Order
	err = json.Unmarshal([]byte(input), &order)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Создаем подключение к NATS
	nc, err := nats.Connect(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Создаем подключение к NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	defer sc.Close()

	// Преобразуем структуру заказа в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Определяем канал и блокировку для асинхронной публикации
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
	guid, err = sc.PublishAsync(subj, orderJSON, acb)
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
	glock.Unlock()
	if guid == "" {
		log.Fatal("Expected non-empty guid to be returned.")
	}
	log.Printf("Published [%s] : '%s' [guid: %s]\n", subj, orderJSON, guid)

	select {
	case <-ch:
		log.Println("Publication acknowledged.")
		break
	case <-time.After(5 * time.Second):
		log.Fatal("timeout")
	}

	// Добавление заказа в базу данных
	err = r.AddToDB(order)
	c.AddToCache(order)
	if err != nil {
		log.Fatalf("Error adding order to database: %v", err)
	}

	fmt.Println("Order successfully added to the database and published to NATS.")
}
