package router

import (
    "github.com/gin-gonic/gin"
    "paralelos/controller"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Definir rutas de usuarios
    r.GET("/users", controller.GetUsers)
    r.GET("/users/:id", controller.GetUserByID)

    return r
}
