package shared

var RabbitMQUrl = "amqp://guest:guest@rabbitmq:5672/"

type CollatzTask struct {
	Number uint64
}
