package repository

import (
    "database/sql"
    "fmt"
    "paralelos/model"
)

// Aquí puedes agregar lógica para interactuar con una base de datos
func GetUsers() ([]model.User, error) {
    db := GetDB()

    // Consulta para obtener los usuarios de la base de datos
    rows, err := db.Query("SELECT id, name, email FROM users")
    if err != nil {
        return nil, fmt.Errorf("Error al obtener usuarios: %v", err)
    }
    defer rows.Close()

    var users []model.User
    for rows.Next() {
        var user model.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, fmt.Errorf("Error al escanear fila: %v", err)
        }
        users = append(users, user)
    }

    return users, nil
}

func GetUserByID(id uint) (model.User, bool, error) {
    db := GetDB()

    // Consulta para obtener un usuario por su ID
    var user model.User
    err := db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return user, false, nil
        }
        return user, false, fmt.Errorf("Error al obtener usuario por ID: %v", err)
    }

    return user, true, nil
}
