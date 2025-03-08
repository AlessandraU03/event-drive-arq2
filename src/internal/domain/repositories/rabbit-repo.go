package repositories

import ( 
	amqp "github.com/rabbitmq/amqp091-go"
)
// Interfaz para interactuar con RabbitMQ
type IRabbitMQ interface {
    Connect() error
    Close() error
    Publish(queueName string, message []byte) error
    Consume(queueName string) (<-chan amqp.Delivery, error)
}
