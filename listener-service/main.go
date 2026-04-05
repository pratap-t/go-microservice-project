package main

import (
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panic(err)
	}

	// watch the event and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
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
