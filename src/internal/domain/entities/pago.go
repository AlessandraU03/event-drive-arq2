package entities

type Pago struct {
    PagoID     int     `json:"pago_id"`
    PedidoID   int     `json:"pedido_id"`
    Estado string  `json:"estado_pago"` // 'pendiente', 'exitoso', 'fallido'
    Monto      float64 `json:"monto"`
    MetodoPago string  `json:"metodo_pago"`
    Fecha      string  `json:"fecha"`
}
