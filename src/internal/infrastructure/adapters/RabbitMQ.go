// RabbitMQ Adapter para publicar y consumir mensajes

package adapters

import (
	"log"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	channel *amqp.Channel
}

func NewRabbitMQAdapter() (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial("amqp://ale:ale05@54.156.170.232/")
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQAdapter{
		channel: channel,
	}, nil
}

// PublishPedido publica un pedido en la cola
func (r *RabbitMQAdapter) PublishPedido(pedido interface{}) error {
	body, err := json.Marshal(pedido)
	if err != nil {
		log.Printf("Error al convertir el pedido a JSON: %v", err)
		return err
	}

	err = r.channel.Publish(
		"",         // Exchange
		"pedidos", // Cola de mensajes
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Error al publicar en RabbitMQ: %v", err)
		return err
	}
	return nil
}

// ConsumePedidos consume los mensajes de la cola
func (r *RabbitMQAdapter) ConsumePedidos() (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		"pedidos", // Cola de mensajes
		"",             // Consumer
		true,           // Auto-ack
		false,          // Exclusive
		false,          // No-local
		false,          // No-wait
		nil,            // Args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
