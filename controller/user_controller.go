package controller

import (
    "net/http"
    "strconv"
    "paralelos/repository"
    "paralelos/viewmodel"
    "github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
    // Capturar tanto los usuarios como el error
    users, err := repository.GetUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
        return
    }

    // Crear una lista de ViewModels
    var userVMs []viewmodel.UserViewModel
    for _, user := range users {
        userVMs = append(userVMs, viewmodel.NewUserViewModel(user))
    }

    c.JSON(http.StatusOK, userVMs)
}


func GetUserByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    // Capturar tanto el usuario, el booleano y el error
    user, found, err := repository.GetUserByID(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuario"})
        return
    }

    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
        return
    }

    // Si el usuario se encontró, devolverlo como un ViewModel
    c.JSON(http.StatusOK, viewmodel.NewUserViewModel(user))
}