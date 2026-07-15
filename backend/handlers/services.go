package handlers

import (
	"encoding/json"
	"net/http"
	"pulseflux/database"
)

type NewServicePayload struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	CheckInterval int    `json:"check_interval"`
}

// ServicesHandler centraliza las peticiones de administración del Dashboard
func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	// Habilitar CORS básico para desarrollo multiplataforma nativo
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		getServices(w, r)
	case "POST":
		createService(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no soportado"})
	}
}

func getServices(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, name, url, check_interval, status FROM services ORDER BY created_at DESC"
	rows, err := database.DB.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, name, url, status string
		var interval int
		if err := rows.Scan(&id, &name, &url, &interval, &status); err == nil {
			list = append(list, map[string]interface{}{
				"id":             id,
				"name":           name,
				"url":            url,
				"check_interval": interval,
				"status":         status,
			})
		}
	}

	json.NewEncoder(w).Encode(list)
}

func createService(w http.ResponseWriter, r *http.Request) {
	var payload NewServicePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.Name == "" || payload.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payload inválido o campos obligatorios vacíos"})
		return
	}

	if payload.CheckInterval <= 0 {
		payload.CheckInterval = 5
	}

	query := "INSERT INTO services (name, url, check_interval) VALUES ($1, $2, $3) RETURNING id"
	var lastInsertID string
	err = database.DB.QueryRow(query, payload.Name, payload.URL, payload.CheckInterval).Scan(&lastInsertID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "La URL ya se encuentra registrada o fallo en BD"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "id": lastInsertID})
}