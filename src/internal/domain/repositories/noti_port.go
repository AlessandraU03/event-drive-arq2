package repositories

import "payment/src/internal/domain/entities"

// Interfaz para el servicio de notificación (Socket.IO)
type INotificationService interface {
    SendPaymentNotification(pago entities.Pago) error
}
