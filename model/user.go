package model

type User struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"` // No queremos que la contraseña se devuelva en la respuesta JSON
}

