package services

import (
	"payment/src/internal/domain/entities"
	"log"
	"encoding/json"
	"net/http"
	"time"
)

type NotificationService struct {}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// Esta función simula la respuesta a un cliente usando Long Polling
func (ns *NotificationService) SendPaymentNotificationLongPolling(w http.ResponseWriter, r *http.Request, pago entities.Pago) {
	if w == nil {
		log.Println("⚠️ Advertencia: Se llamó a SendPaymentNotification sin un ResponseWriter válido")
		return
	}

	timeout := time.NewTimer(30 * time.Second) 
	done := make(chan bool)

	// Simula que el pedido es procesado después de un tiempo (puedes reemplazar esto con la lógica real)
	go func() {
		time.Sleep(5 * time.Second) // Simulación del tiempo que tarda en procesarse el pedido
		done <- true                // Cuando el pedido se haya procesado, se envía una señal al canal
	}()

	select {
	case <-done:
		// Responde al cliente con la información del pedido procesado
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"pedido_id": pago.PedidoID,
			"estado":    pago.Estado,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Error al enviar la notificación:", err)
		}
	case <-timeout.C:
		// Si el timeout se alcanza antes de que se procese el pedido, respondemos con un error
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("Timeout alcanzado"))
	}
}
