package controller

import (
	"net/http"
	"paralelos/repository"
	"paralelos/viewmodel"
	"strconv"
	"sync"

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


func GetUserNames(c *gin.Context) {
    db := repository.GetDB()

    // Realiza la consulta para obtener los nombres
    rows, err := db.Query("SELECT name FROM users")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los nombres"})
        return
    }
    defer rows.Close()

    var names []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar los nombres"})
            return
        }
        names = append(names, name)
    }

    c.JSON(http.StatusOK, names)
}

func GetUserEmails(c *gin.Context) {
    db := repository.GetDB()

    // Realiza la consulta para obtener los correos electrónicos
    rows, err := db.Query("SELECT email FROM users")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los correos"})
        return
    }
    defer rows.Close()

    var emails []string
    for rows.Next() {
        var email string
        if err := rows.Scan(&email); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar los correos"})
            return
        }
        emails = append(emails, email)
    }

    c.JSON(http.StatusOK, emails)
}




// Estructura para devolver el nombre y correo concatenados
type UserConcat struct {
    NameAndEmail string `json:"name_and_email"`
}

// Handler para consultar name y email en paralelo y concatenarlos
func GetUserNamesAndEmails(c *gin.Context) {
    db := repository.GetDB()

    var wg sync.WaitGroup
    var names []string
    var emails []string
    var nameErr, emailErr error

    wg.Add(2)

    // Goroutine para obtener los nombres
    go func() {
        defer wg.Done()
        rows, err := db.Query("SELECT name FROM users")
        if err != nil {
            nameErr = err
            return
        }
        defer rows.Close()

        for rows.Next() {
            var name string
            if err := rows.Scan(&name); err != nil {
                nameErr = err
                return
            }
            names = append(names, name)
        }
    }()

    // Goroutine para obtener los correos electrónicos
    go func() {
        defer wg.Done()
        rows, err := db.Query("SELECT email FROM users")
        if err != nil {
            emailErr = err
            return
        }
        defer rows.Close()

        for rows.Next() {
            var email string
            if err := rows.Scan(&email); err != nil {
                emailErr = err
                return
            }
            emails = append(emails, email)
        }
    }()

    // Esperar a que ambas goroutines terminen
    wg.Wait()

    if nameErr != nil || emailErr != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los datos"})
        return
    }

    // Combinar los nombres y correos electrónicos
    if len(names) != len(emails) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "El número de nombres y correos no coinciden"})
        return
    }

    var usersConcat []UserConcat
    for i := 0; i < len(names); i++ {
        usersConcat = append(usersConcat, UserConcat{
            NameAndEmail: names[i] + " - " + emails[i],
        })
    }

    c.JSON(http.StatusOK, usersConcat)
}