package router

import (
    "github.com/gin-gonic/gin"
    "paralelos/controller"
    "paralelos/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.Use(middleware.CORSMiddleware())

    // Rutas públicas
    r.POST("/register", controller.Register)
    r.POST("/login", controller.Login)

    // Rutas protegidas
    authorized := r.Group("/")
    authorized.Use(middleware.AuthMiddleware())
    {
        authorized.GET("/users", controller.GetUsers)
        authorized.GET("/users/:id", controller.GetUserByID)
    }

    return r
}
