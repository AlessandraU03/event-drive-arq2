// infrastructure/routes/routes.go

package routes

import (
	"log"
	"payment/src/internal/application/services"
	"payment/src/internal/infrastructure/adapters"
	"payment/src/internal/infrastructure/controllers"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

// SetupRouter configura las rutas de la API
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Crear el adaptador RabbitMQ
	rabbitAdapter, err := adapters.NewRabbitMQAdapter()
	if err != nil {
		log.Fatal("Error al crear el adaptador RabbitMQ: ", err)
	}

	// Crear el servidor Socket.io
	socketServer := socketio.NewServer(nil)

	// Crear el servicio de notificación (que utiliza Socket.io)
	notificationService := services.NewNotificationService(socketServer)

	// Crear el controlador de procesamiento de pagos
	paymentProcessorController := controllers.NewPaymentProcessorController(rabbitAdapter, socketServer, notificationService)

	// Ruta para procesar los pagos (aunque en este caso el consumidor de RabbitMQ ya está escuchando)
	// Esta ruta se puede agregar si se quiere hacer una validación o integración adicional
	router.POST("/procesar_pago", func(c *gin.Context) {
		paymentProcessorController.ConsumePedido()
	})

	// Configuración del WebSocket para notificaciones en tiempo real
	router.GET("/socket.io/", gin.WrapH(socketServer))

	// Rutas para escuchar el procesamiento de pedidos y pagos
	go paymentProcessorController.ConsumePedido()

	return router
}
