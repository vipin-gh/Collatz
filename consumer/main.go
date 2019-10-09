package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/icarus3/Collatz/shared"
	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	time.Sleep(20 * time.Second)

	conn, err := amqp.Dial(shared.RabbitMQUrl)
	handleError(err, "Can't connect to AMQ")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("collatz", true, false, false, false, nil)
	handleError(err, "Could not declare `collatz` queue")

	err = amqpChannel.Qos(1,0,false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)
			
			collatzTask := &shared.CollatzTask{}

			err := json.Unmarshal(d.Body, collatzTask)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			steps := 0
			num := collatzTask.Number
			
			for num != 1 {
				if num % 2 == 0 {
					num = num / 2
				} else {
					num = ( num * 3 ) + 1
				}
				steps = steps + 1
			}

			log.Printf("Number of steps took for %d are %d", collatzTask.Number, steps)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	<-stopChan


}

