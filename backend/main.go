package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	
	"pulseflux/database" // Importamos nuestro paquete interno de base de datos
)

//go:embed dist/*
var frontendFiles embed.FS

func main() {
	fmt.Println("=== [PulseFlux Core Engine] Iniciando Servicio ===")

	// 1. Inicializar Base de Datos (Módulo de Persistencia)
	// Los scripts run.bash / run.bat aseguran que las variables de entorno existan
	database.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Configurar endpoints de la API (Continuará en el siguiente módulo)
	http.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "success", "database": "connected"}`))
	})

	// 3. Servir frontend embebido
	publicFS, err := fs.Sub(frontendFiles, "dist")
	if err != nil {
		log.Fatalf("Error al cargar la interfaz embebida: %v", err)
	}
	http.Handle("/", http.FileServer(http.FS(publicFS)))

	fmt.Printf("Servidor desplegado exitosamente en: http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Fallo crítico en el servidor: %v", err)
	}
}