package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver oficial puro de Go para PostgreSQL
)

var DB *sql.DB

// InitDB inicializa un pool de conexiones optimizado y seguro para hilos concurrente
func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")

	// Cadena de conexión estándar (DSN)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, user, name)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error crítico al abrir la conexión con PostgreSQL: %v", err)
	}

	// Configuración del Pool de Conexiones a nivel Senior
	DB.SetMaxOpenConns(25)                 // Límite de conexiones simultáneas abiertas
	DB.SetMaxIdleConns(25)                 // Mantener conexiones inactivas calientes para reuso rápido
	DB.SetConnMaxLifetime(5 * time.Minute) // Evitar fugas de memoria o conexiones corruptas viejas

	// Verificar conexión real contra el servidor
	err = DB.Ping()
	if err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL de forma activa: %v", err)
	}

	fmt.Println("[+] Pool de conexiones a PostgreSQL inicializado con éxito.")
}