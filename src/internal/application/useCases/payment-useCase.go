package useCases

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
    pago.Estado = "exitoso"

    err := uc.PagoRepository.SavePago(*pago)
    if err != nil {
        return err
    }

    return nil
}
