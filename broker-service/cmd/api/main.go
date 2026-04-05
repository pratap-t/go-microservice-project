package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	// try to connect rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp091.Connection, error) {
	var counts int64
	var backoff = time.Second * 1
	var connection *amqp091.Connection

	// don't continue until rabbitmq is ready
	for {
		c, err := amqp091.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ!")
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}
