package repositories

import "payment/src/internal/domain/entities"

type IPedido interface {
	Save(pedido *entities.Pedido) error

}