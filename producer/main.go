package main

import (
	"encoding/json"
	"log"
	"math"
	"fmt"
	"time"

	"github.com/icarus3/Collatz/shared"		
	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func retryConn() (*amqp.Connection, error) {
	retry := 10
	for retry > 0 {
		conn, err := amqp.Dial(shared.RabbitMQUrl)
		if err == nil {
			return conn, err
		}

		retry = retry - 1
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("connection timeout")
}

func main() {
	conn, err := retryConn()
	handleError(err, "Can't connect to AMQ")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("collatz", true, false, false, false, nil)
	handleError(err, "Could not declare `collatz` queue")

	for num := uint64(1); num < math.MaxUint64; num++ {
		collatzTask := shared.CollatzTask{Number: num}
		body, err := json.Marshal(collatzTask)

		if err != nil {
			handleError(err, "Error encoding json")
		}

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing {
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body: body,	
		})

		if err != nil {
			log.Fatalf("Error publishing message: %s", err)
		}

		log.Printf("CollatzTask: %d", collatzTask.Number)
	}
}