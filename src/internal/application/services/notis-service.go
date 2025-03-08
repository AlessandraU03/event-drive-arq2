// infrastructure/services/notification_service.go
package services

import (
    "payment/src/internal/domain/entities"
    "github.com/googollee/go-socket.io"
)

type NotificationService struct {
    Server *socketio.Server
}

func NewNotificationService(server *socketio.Server) *NotificationService {
    return &NotificationService{Server: server}
}

func (ns *NotificationService) SendPaymentNotification(pago entities.Pago) error {
    ns.Server.BroadcastToNamespace("", "/", "payment_status", map[string]interface{}{
        "pedido_id": pago.PedidoID,
        "estado":    pago.Estado,
    })
    return nil
}
