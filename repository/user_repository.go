package repository

import (
    "database/sql"
    "fmt"
    "paralelos/model"
    "golang.org/x/crypto/bcrypt"
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

// Registrar un nuevo usuario
func RegisterUser(name, email, password string) error {
    db := GetDB()

    // Verificar si el usuario ya existe
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
    if err != nil {
        return fmt.Errorf("Error al verificar usuario: %v", err)
    }

    if exists {
        return fmt.Errorf("El usuario con el email ya existe")
    }

    // Generar un hash de la contraseña
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("Error al generar hash de la contraseña: %v", err)
    }

    // Insertar el nuevo usuario en la base de datos
    _, err = db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", name, email, string(hashedPassword))
    if err != nil {
        return fmt.Errorf("Error al registrar usuario: %v", err)
    }

    return nil
}


func AuthenticateUser(email, password string) (model.User, bool, error) {
    db := GetDB()

    // Buscar el usuario por email
    var user model.User
    err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            return user, false, nil // Usuario no encontrado
        }
        return user, false, fmt.Errorf("Error al autenticar usuario: %v", err)
    }

    // Comparar la contraseña proporcionada con el hash almacenado en la base de datos
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return user, false, nil // Contraseña incorrecta
    }

    return user, true, nil // Usuario autenticado correctamente
}