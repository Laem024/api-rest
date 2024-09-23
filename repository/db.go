package repository

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
    var err error
    // Cadena de conexión
    connStr := "user=postgres password=root dbname=paralelo sslmode=disable port=5432"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("Error al conectar a la base de datos: %v", err)
    }

    // Verificar que la conexión sea correcta
    err = db.Ping()
    if err != nil {
        return fmt.Errorf("Error al verificar la base de datos: %v", err)
    }

    fmt.Println("Conectado exitosamente a PostgreSQL")
    return nil
}

func GetDB() *sql.DB {
    return db
}
