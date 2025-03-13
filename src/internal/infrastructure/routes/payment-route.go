package routes

import (
	"log"
	"payment/src/internal/application/services"
	"payment/src/internal/infrastructure/adapters"
	"payment/src/internal/infrastructure/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, 
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Crear el adaptador RabbitMQ
	rabbitAdapter, err := adapters.NewRabbitMQAdapter()
	if err != nil {
		log.Fatal("Error al crear el adaptador RabbitMQ: ", err)
	}

	// Crear el servicio de notificación
	notificationService := services.NewNotificationService()

	// Crear el controlador de procesamiento de pagos
	paymentProcessorController := controllers.NewPaymentProcessorController(rabbitAdapter, notificationService)

	// Ruta para procesar pagos
	router.GET("/procesar_pago", func(c *gin.Context) {
		pedidoID := c.DefaultQuery("pedido_id", "")
		if pedidoID == "" {
			c.JSON(400, gin.H{"error": "El parámetro 'pedido_id' es obligatorio"})
			return
		}

		log.Println("✅ Procesando pedido con ID:", pedidoID)
		// Ejecutar el controlador para consumir procesados
		paymentProcessorController.ConsumeProcesados(c.Writer, c.Request)

	})

	// Iniciar el consumidor de RabbitMQ en segundo plano
	go paymentProcessorController.ProcessarPedidosEnSegundoPlano()

	return router
}
