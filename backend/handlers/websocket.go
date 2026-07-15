package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

// ClientHub gestiona todos los túneles WebSockets abiertos simultáneamente
type ClientHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.Mutex
}

var Hub = ClientHub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan []byte),
}

// StartHubEngine inicia un hilo asíncrono para despachar mensajes a los clientes conectados
func StartHubEngine() {
	go func() {
		for {
			message := <-Hub.broadcast
			Hub.mutex.Lock()
			for client := range Hub.clients {
				err := websocket.Message.Send(client, string(message))
				if err != nil {
					log.Printf("[-] Error al enviar trama por WebSocket, cerrando canal: %v", err)
					client.Close()
					delete(Hub.clients, client)
				}
			}
			Hub.mutex.Unlock()
		}
	}()
}

// BroadcastMessage es la función pública que el motor de monitoreo invocará para enviar datos dinámicos
func BroadcastMessage(jsonPayload string) {
	Hub.broadcast <- []byte(jsonPayload)
}

// WebSocketHandler gestiona el Upgrade HTTP y acopla el cliente al Hub
func WebSocketHandler(ws *websocket.Conn) {
	Hub.mutex.Lock()
	Hub.clients[ws] = true
	Hub.mutex.Unlock()

	fmt.Printf("[+] Cliente web acoplado al túnel WebSocket. Total activos: %d\n", len(Hub.clients))

	// Mantener la conexión abierta escuchando tramas vacías (Heartbeat) para evitar cierres prematuros
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			Hub.mutex.Lock()
			delete(Hub.clients, ws)
			Hub.mutex.Unlock()
			fmt.Println("[-] Cliente web desconectado del túnel WebSocket.")
			break
		}
	}
}