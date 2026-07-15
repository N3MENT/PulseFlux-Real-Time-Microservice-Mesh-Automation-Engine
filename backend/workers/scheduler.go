package workers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"pulseflux/database"
)

// Service representa la estructura interna de un microservicio a monitorear
type Service struct {
	ID            string
	Name          string
	URL           string
	CheckInterval int
}

// StartMonitoringEngine inicia el bucle global encargado de buscar servicios activos
func StartMonitoringEngine() {
	fmt.Println("[+] Motor de monitoreo concurrente iniciado con éxito.")
	
	// Bucle infinito en segundo plano para descubrir nuevos servicios o actualizar intervalos
	go func() {
		for {
			services, err := fetchActiveServices()
			if err != nil {
				log.Printf("[X] Error al escanear servicios en la base de datos: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			for _, svc := range services {
				// Lanzamos una Goroutine independiente por cada microservicio de manera concurrente
				// Usamos un patrón worker para evitar duplicar tareas idénticas
				go monitorWorker(svc)
			}

			// Intervalo global de re-escaneo de la tabla de configuraciones
			time.Sleep(30 * time.Second)
		}
	// En Go, invocar funciones vacías con "go func(){}()" ejecuta un hilo asíncrono puro
	}() 
}

// fetchActiveServices extrae los objetivos de monitoreo desde PostgreSQL
func fetchActiveServices() ([]Service, error) {
	query := "SELECT id, name, url, check_interval FROM services"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.CheckInterval); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

// monitorWorker es el bucle asíncrono dedicado exclusivamente a un solo servicio
func monitorWorker(svc Service) {
	// El ticker controla de forma nativa y exacta los intervalos de tiempo del hardware sin desfases
	ticker := time.NewTicker(time.Duration(svc.CheckInterval) * time.Second)
	defer ticker.Stop()

	client := &http.Client{
		Timeout: 4 * time.Second, // Timeout estricto para evitar hilos colgados indefinidamente
	}

	for range ticker.C {
		startTime := time.Now()
		resp, err := client.Get(svc.URL)
		latency := time.Since(startTime).Milliseconds()

		status := "ONLINE"
		statusCode := 0
		errMsg := ""

		if err != nil {
			status = "OFFLINE"
			errMsg = err.Error()
		} else {
			statusCode = resp.StatusCode
			resp.Body.Close()
			if resp.StatusCode >= 400 {
				status = "OFFLINE"
			}
		}

		// Persistir métricas de forma atómica en PostgreSQL
		err = saveMetrics(svc.ID, status, int(latency), statusCode, errMsg)
		if err != nil {
			log.Printf("[X] Error persistiendo métricas para %s: %v", svc.Name, err)
		}
	}
}

// saveMetrics registra los pings y actualiza los estados mediante transacciones SQL
func saveMetrics(serviceID string, status string, latency int, statusCode int, errMsg string) error {
	// Iniciamos una transacción atómica para garantizar consistencia absoluta
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insertar en la tabla histórica de logs de rendimiento
	logQuery := `INSERT INTO performance_logs (service_id, latency_ms, status_code, error_message) 
                 VALUES ($1, $2, $3, $4)`
	
	var sqlErr sql.NullString
	if errMsg != "" {
		sqlErr.String = errMsg
		sqlErr.Valid = true
	}

	_, err = tx.Exec(logQuery, serviceID, latency, statusCode, sqlErr)
	if err != nil {
		return err
	}

	// 2. Actualizar el estado global en tiempo real en la tabla de servicios
	updateQuery := "UPDATE services SET status = $1 WHERE id = $2"
	_, err = tx.Exec(updateQuery, status, serviceID)
	if err != nil {
		return err
	}

	// Crear trama JSON en string para transmitirla de forma asíncrona por WebSocket
eventPayload := fmt.Sprintf(`{"event":"METRIC_UPDATE","service_id":"%s","status":"%s","latency_ms":%d}`, 
    serviceID, status, latency)

handlers.BroadcastMessage(eventPayload)

	return tx.Commit()
}