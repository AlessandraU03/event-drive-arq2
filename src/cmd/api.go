package cmd


import (
    "payment/src/internal/infrastructure/routes"
    "log"
 
)

func Api() {
    router := routes.SetupRouter() // 🚀 Usa el router de SetupRouter()

    // Iniciar el servidor HTTP en el puerto 8081
    log.Println("✅ Servidor corriendo en http://localhost:8081")
    log.Fatal(router.Run(":8081"))
}

