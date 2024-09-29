package controller

import (
    "time"
    "net/http"
    "paralelos/repository"
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
    var req RegisterRequest

    // Validar los datos del registro
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    // Registrar al usuario en la base de datos
    err := repository.RegisterUser(req.Name, req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado con éxito"})
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

var JwtSecret = []byte("nsafjnsajfsajn12421421n421jn4kj1rjqnfiqjr021i4n21i4n1or21oinr") // Ahora esta variable es exportada (pública)

func generateJWT(userEmail string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": userEmail,
        "exp":   time.Now().Add(time.Hour * 72).Unix(), // Token válido por 72 horas
    })

    // Firmar el token con la clave secreta
    tokenString, err := token.SignedString(JwtSecret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func Login(c *gin.Context) {
    var req LoginRequest

    // Validar el request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
        return
    }

    // Autenticar al usuario
    user, authenticated, err := repository.AuthenticateUser(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al autenticar usuario"})
        return
    }

    if !authenticated {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña incorrectos"})
        return
    }

    // Generar el token JWT
    token, err := generateJWT(user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Autenticación exitosa",
        "token":   token,
    })
}
