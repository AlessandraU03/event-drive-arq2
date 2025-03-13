package adapters

import (
	"payment/src/core"
	"payment/src/internal/domain/entities"
	"fmt"
	"log"
)

type MySQLPedidos struct {
	conn *core.Conn_MySQL
}

// Constructor para MySQLPedidos
func NewMySQLPedidos() *MySQLPedidos {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQLPedidos{conn: conn}
}

// Guardar un nuevo pedido en la base de datos
func (mysql *MySQLPedidos) Save(pedido *entities.Pedido) error {
	query := `INSERT INTO pedidos (producto_id, cliente, metodo_pago, monto_pagado, direccion, estado) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	result, err := mysql.conn.ExecutePreparedQuery(query, pedido.ProductoID, pedido.Cliente, pedido.MetodoPago, pedido.Monto, pedido.Direccion, pedido.Estado)
	if err != nil {
		return fmt.Errorf("error al guardar el pedido: %w", err)
	}

	// Obtener el ID generado
	pedidoID, _ := result.LastInsertId()
	pedido.ID = int(pedidoID)  // Asignar el ID generado al pedido
	log.Printf("[MySQLPedidos] - Pedido guardado con ID: %d", pedido.ID)
	return nil
}
