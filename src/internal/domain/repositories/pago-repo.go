package repositories

import "payment/src/internal/domain/entities"

// Interfaz para el repositorio de pagos
type IPagoRepository interface {
    SavePago(pago entities.Pago) error
    UpdatePagoEstado(pagoID int, estado string) error
    GetPagoByPedidoID(pedidoID int) (*entities.Pago, error)
}
