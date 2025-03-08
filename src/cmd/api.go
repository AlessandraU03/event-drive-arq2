package cmd


import (
    "payment/src/internal/infrastructure/routes"
    "log"
	"github.com/gin-gonic/gin"
)

func Api() {
    // Configurar rutas
    router := gin.Default()
    routes.SetupRouter()


    // Iniciar el servidor HTTP
    log.Fatal(router.Run(":8081"))
}
