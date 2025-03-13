package controllers

import (
	"encoding/json"
	"log"
	"payment/src/internal/application/services"
	"payment/src/internal/domain/entities"
	"payment/src/internal/infrastructure/adapters"
	"net/http"
	"time"
)

type PaymentProcessorController struct {
	rabbitAdapter      *adapters.RabbitMQAdapter
	notificationService *services.NotificationService
}

func NewPaymentProcessorController(
	rabbitAdapter *adapters.RabbitMQAdapter,
	notificationService *services.NotificationService,
) *PaymentProcessorController {
	return &PaymentProcessorController{
		rabbitAdapter:      rabbitAdapter,
		notificationService: notificationService,
	}
}

// ConsumePedido procesa los pedidos de RabbitMQ y los envía a través de Long Polling
func (ppc *PaymentProcessorController) ConsumePedido(w http.ResponseWriter, r *http.Request) {
    msgs, err := ppc.rabbitAdapter.ConsumePedidos()
    if err != nil {
        log.Fatalf("❌ Error al consumir pedidos: %v", err)
    }

    // Configurar Long Polling
    w.Header().Set("Content-Type", "application/json")
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
    }

    // Canal de notificación
    notifyChan := make(chan entities.Pedido)

    // Enviar la notificación de pago cuando un pedido sea procesado
    go func() {
        for msg := range msgs {
            var pedido entities.Pedido

            log.Println("📥 Mensaje recibido:", string(msg.Body))

            unmarshalErr := json.Unmarshal(msg.Body, &pedido)
            if unmarshalErr != nil {
                log.Printf("❌ Error al procesar el pedido: %v", unmarshalErr)
                msg.Nack(false, false)
                continue
            }

            log.Printf("✅ Pedido procesado: %+v", pedido)
            msg.Ack(false)

            // Enviar el pedido procesado al canal de notificación
            notifyChan <- pedido
        }
    }()

    // Long Polling: Mantener la conexión abierta hasta que llegue una notificación
    select {
    case pedido := <-notifyChan:
        if err := json.NewEncoder(w).Encode(pedido); err != nil {
            log.Printf("❌ Error al enviar el pedido por Long Polling: %v", err)
        }
        flusher.Flush()
    case <-time.After(30 * time.Second): // Timeout si no hay mensaje
        http.Error(w, "Timeout alcanzado", http.StatusRequestTimeout)
    }
}

func (ppc *PaymentProcessorController) ConsumeProcesados(w http.ResponseWriter, r *http.Request) {
	// Consumir los mensajes de la cola "procesados"
	msgs, err := ppc.rabbitAdapter.ConsumeProcesados()
	if err != nil {
		log.Fatalf("❌ Error al consumir la cola 'procesados': %v", err)
	}

	// Configuración de Long Polling
	w.Header().Set("Content-Type", "application/json")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Leer mensajes de la cola "procesados" y enviarlos al cliente
	for msg := range msgs {
		var pedido entities.Pedido

		// Deserializar el mensaje recibido de RabbitMQ
		if err := json.Unmarshal(msg.Body, &pedido); err != nil {
			log.Printf("❌ Error al procesar el pedido: %v", err)
			msg.Nack(false, false)
			continue
		}

		// Confirmar que el mensaje fue procesado
		log.Printf("✅ Pedido procesado: %+v", pedido)
		msg.Ack(false)

		// Enviar el pedido a través de Long Polling
		if err := json.NewEncoder(w).Encode(pedido); err != nil {
			log.Printf("❌ Error al enviar el pedido: %v", err)
			break
		}

		// Forzar el envío de la respuesta al cliente
		flusher.Flush()

		// Simular un tiempo de espera antes de procesar el siguiente mensaje
		time.Sleep(2 * time.Second)
	}
}

func (ppc *PaymentProcessorController) ProcessarPedidosEnSegundoPlano() {
	msgs, err := ppc.rabbitAdapter.ConsumePedidos()
	if err != nil {
		log.Fatalf("❌ Error al consumir pedidos: %v", err)
	}

	for msg := range msgs {
		var pedido entities.Pedido
		log.Println("📥 Mensaje recibido:", string(msg.Body))

		unmarshalErr := json.Unmarshal(msg.Body, &pedido)
		if unmarshalErr != nil {
			log.Printf("❌ Error al procesar el pedido: %v", unmarshalErr)

			msg.Nack(false, false)
			continue
		}

		log.Printf("✅ Pedido procesado: %+v", pedido)
		msg.Ack(false)

		err := ppc.rabbitAdapter.PublishToQueue("procesados", msg.Body)
		if err != nil {
			log.Printf("⚠️ Error al reenviar mensaje a la cola 'procesados': %v", err)
			msg.Nack(false, true) 
			continue
		}

		time.Sleep(2 * time.Second)
	}
}
