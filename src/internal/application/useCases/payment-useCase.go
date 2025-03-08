// application/usecases/process_payment.go
package usecases

import (
	"payment/src/internal/application/services"
	"payment/src/internal/domain/repositories"
	"payment/src/internal/domain/entities"
)



type ProcessPayment struct {
    PagoRepository     repositories.IPagoRepository
    NotificationService services.NotificationService
}

func NewProcessPayment(pagoRepo repositories.IPagoRepository, notificationService services.NotificationService) *ProcessPayment {
    return &ProcessPayment{
        PagoRepository:     pagoRepo,
        NotificationService: notificationService,
    }
}

func (uc *ProcessPayment) Execute(pago *entities.Pago) error {
    // Lógica de validación del pago (simulada)
    pago.Estado = "exitoso"

    // Guardar el pago en la base de datos
    err := uc.PagoRepository.SavePago(*pago)
    if err != nil {
        return err
    }

    // Emitir notificación al frontend sobre el estado del pago
    err = uc.NotificationService.SendPaymentNotification(*pago)
    if err != nil {
        return err
    }

    return nil
}
