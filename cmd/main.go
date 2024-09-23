package main

import (
    "log"
    "paralelos/repository"
    "paralelos/router"
)

func main() {
    // Inicializar la conexión a la base de datos
    if err := repository.InitDB(); err != nil {
        log.Fatalf("No se pudo conectar a la base de datos: %v", err)
    }

    // Configurar el router de Gin
    r := router.SetupRouter()

    // Ejecutar la aplicación en el puerto 8080
    r.Run(":8090")
}
