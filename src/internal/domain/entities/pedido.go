package entities

type Pedido struct {
	ID          int        `json:"pedido_id"`
	ProductoID  int        `json:"producto_id"` // Producto asociado al pedido
	Cliente     string     `json:"cliente"`     // Nombre del cliente
	Direccion   string     `json:"direccion"`
	MetodoPago  string     `json:"metodo_pago"` // 'Tarjeta', 'Efectivo'
	Monto       float64    `json:"monto_pagado"`
	Estado      string     `json:"estado"` // Pendiente, En preparación, Listo
}

func NewPedido(pedidoID int, productoID int, cliente string, direccion string, metodoPago string, monto float64, estado string) *Pedido {
	return &Pedido{
		ID:         pedidoID,
		ProductoID: productoID,
		Cliente:    cliente,
		Direccion:  direccion,
		MetodoPago: metodoPago,
		Monto:      monto,
		Estado:     estado,
	}
}
