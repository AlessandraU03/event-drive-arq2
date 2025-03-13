package adapters

import (

	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQAdapter() (*RabbitMQAdapter, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://ale:ale05@54.156.170.232/")
	if err != nil {
		return nil, fmt.Errorf("error al conectar a RabbitMQ: %w", err)
	}

	// Crear un canal
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error al abrir canal: %w", err)
	}

	return &RabbitMQAdapter{
		conn:    conn,
		channel: channel,
	}, nil
}

// ConsumePedidos consume los mensajes de la cola "pedidos"
func (r *RabbitMQAdapter) ConsumePedidos() (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(
		"pedidos", 
		true,      
		false,     
		false,     
		false,     
		nil,       
	)
	if err != nil {
		return nil, fmt.Errorf("error al declarar la cola 'pedidos': %w", err)
	}

	// Consumir los mensajes de la cola "pedidos"
	msgs, err := r.channel.Consume(
		"pedidos", 
		"",        
		false,     
		false,    
		false,    
		false,     
		nil,       
	)
	if err != nil {
		return nil, fmt.Errorf("error al consumir mensajes de la cola 'pedidos': %w", err)
	}

	return msgs, nil
}

// ConsumeProcesados consume los mensajes de la cola "procesados"
func (r *RabbitMQAdapter) ConsumeProcesados() (<-chan amqp.Delivery, error) {
	
	_, err := r.channel.QueueDeclare(
		"procesados", 
		true,         
		false,        
		false,        
		false,        
		nil,         
	)
	if err != nil {
		return nil, fmt.Errorf("error al declarar la cola 'procesados': %w", err)
	}

	// Consumir los mensajes de la cola "procesados"
	msgs, err := r.channel.Consume(
		"procesados", 
		"",           
		false,        
		false,       
		false,        
		false,        
		nil,          
	)
	if err != nil {
		return nil, fmt.Errorf("error al consumir mensajes de la cola 'procesados': %w", err)
	}

	return msgs, nil
}


// PublishToQueue publica un mensaje en la cola
func (r *RabbitMQAdapter) PublishToQueue(queueName string, body []byte) error {
	_, err := r.channel.QueueDeclare(
		queueName, 
		true,      
		false,    
		false,     
		false,     
		nil,       
	)
	if err != nil {
		return fmt.Errorf("error al declarar la cola '%s': %w", queueName, err)
	}

	
	err = r.channel.Publish(
		"",        
		queueName, 
		false,     
		false,     
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error al publicar el mensaje en la cola '%s': %w", queueName, err)
	}

	log.Printf("✅ Mensaje enviado a la cola '%s'", queueName)
	return nil
}

func (r *RabbitMQAdapter) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("Conexión cerrada con RabbitMQ")
}
