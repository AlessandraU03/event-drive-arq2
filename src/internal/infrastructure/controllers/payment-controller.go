// API 2 - Consumidor que valida el pago y notifica al cliente

package controllers

import (
	"log"
	"payment/src/internal/application/services"
	"payment/src/internal/infrastructure/adapters"
	"payment/src/internal/domain/entities"
	"encoding/json"
	"github.com/googollee/go-socket.io"
	
)

type PaymentProcessorController struct {
	rabbitAdapter      *adapters.RabbitMQAdapter
	socketServer       *socketio.Server
	notificationService *services.NotificationService
}

func NewPaymentProcessorController(rabbitAdapter *adapters.RabbitMQAdapter, socketServer *socketio.Server, notificationService *services.NotificationService) *PaymentProcessorController {
	return &PaymentProcessorController{
		rabbitAdapter:      rabbitAdapter,
		socketServer:       socketServer,
		notificationService: notificationService,
	}
}

// ConsumePedido consume el pedido de RabbitMQ
func (ppc *PaymentProcessorController) ConsumePedido() {
	msgs, err := ppc.rabbitAdapter.ConsumePedidos()
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		// El mensaje recibido desde RabbitMQ contiene los detalles del pedido
		var pedido struct {
			PedidoID    int     `json:"pedido_id"`
			ClienteID   int     `json:"cliente_id"`
			Total       float64 `json:"total"`
			Estado      string  `json:"estado"`
			Direccion   string  `json:"direccion"`
			MetodoPago  string  `json:"metodo_pago"`
		}

		err := json.Unmarshal(msg.Body, &pedido)
		if err != nil {
			log.Printf("Error al procesar el pedido: %v", err)
			continue
		}

		// Validar el pago (esto es un ejemplo, puedes validarlo de la manera que necesites)
		if pedido.Total > 0 {
			// Aquí iría la validación del pago...

			// Notificar al cliente vía Socket.io
			pago := entities.Pago{
				PedidoID: pedido.PedidoID,
				Estado:   pedido.Estado,
			}
			err := ppc.notificationService.SendPaymentNotification(pago)
			if err != nil {
				log.Printf("Error al enviar la notificación: %v", err)
			}
		}

		// Confirmamos que el mensaje ha sido procesado correctamente
		msg.Ack(false)
	}
}

