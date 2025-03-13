package useCases

import (
	"payment/src/internal/domain/entities"
	"payment/src/internal/domain/repositories"
)



type CreateOrderUseCase struct {
	orderRepository repositories.IPedido
} 

func NewCreateOrderUseCase(orderRepository repositories.IPedido) *CreateOrderUseCase {
	return &CreateOrderUseCase{orderRepository: orderRepository}
}

func (useCase *CreateOrderUseCase) Execute(pedido *entities.Pedido) {
	useCase.orderRepository.Save(pedido)
}