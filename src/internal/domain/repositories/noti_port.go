package repositories

import "payment/src/internal/domain/entities"

type INotificationService interface {
    SendPaymentNotification(pago entities.Pago) error
}
