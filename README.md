# PulseFlux-Real-Time-Microservice-Mesh-Automation-Engine
malla de servidores y automatización en tiempo real 
# PulseFlux – Real-Time Microservice Mesh & Automation Engine

PulseFlux es una plataforma de ingeniería de software diseñada para el monitoreo concurrente de servicios, análisis de latencia en tiempo real y automatización de flujos de trabajo mediante WebSockets y PostgreSQL.

---

## REQUISITOS DEL SISTEMA (Pre-instalación)

Para ejecutar esta aplicación, el sistema anfitrión **solo debe contar con**:
1. **PostgreSQL** (Versión 14 o superior) corriendo en el puerto 5432.
2. **Go (Golang)** instalado (solo en la máquina de desarrollo/compilación).

*Nota de Arquitectura: El cliente final no necesita instalar Node.js, Vite, Nginx ni ninguna otra herramienta de terceros. El frontend se compila a HTML/JS estático y el compilador de Go lo introduce directamente dentro del binario ejecutable final (.exe en Windows / ejecutable nativo en Linux).*

---

## ARQUITECTURA Y DETALLE DE MODULARIZACIÓN

El proyecto está rígidamente estructurado bajo el principio de separación de responsabilidades en los siguientes módulos:

### 1. Módulo de Configuración (`backend/config/`)
* **Propósito:** Centraliza la lectura de variables de entorno y parámetros del sistema (credenciales de base de datos, puertos de escucha, intervalos de muestreo). 
* **Por qué existe:** Evita tener valores fijos (*hardcoded*) en el código y asegura que la aplicación pueda mutar su comportamiento entre desarrollo y producción sin alterar la lógica.

### 2. Módulo de Persistencia y Modelos (`backend/database/`)
* **Propósito:** Gestiona el ciclo de vida de las conexiones (Connection Pool) hacia PostgreSQL, ejecuta las migraciones iniciales de las tablas y aloja los modelos de datos.
* **Características Avanzadas:** Diseñado para manejar queries de alta frecuencia y configurar escuchas pasivas mediante `LISTEN/NOTIFY` de PostgreSQL para reaccionar a eventos de la base de datos de forma asíncrona.

### 3. Módulo de Controladores y Enrutamiento (`backend/handlers/`)
* **Propósito:** Define los puntos de entrada (Endpoints) de la API REST (`GET /api/services`, `POST /api/services`) y gestiona el ciclo de vida de las conexiones **WebSocket**.
* **Por qué existe:** Separa los protocolos de comunicación de la lógica de negocio. Transforma las peticiones HTTP externas en estructuras de datos legibles por el sistema y realiza el *upgrade* de conexiones HTTP a WebSockets Full-Duplex.

### 4. Módulo de Monitoreo Concurrente (`backend/workers/`)
* **Propósito:** Es el motor del sistema. Ejecuta "Goroutines" (hilos ligeros del sistema) que corren de forma indefinida en segundo plano, realizando pings e inspecciones de red (HTTP/TCP) a los servicios registrados de manera asíncrona sin bloquear la API.

### 5. Interfaz de Usuario (`frontend/`)
* **Propósito:** SPA (Single Page Application) que consume la API y los WebSockets.
* **Diseño y Estética:** Diseñado bajo una estética Cyberpunk con paneles de vidrio difuminado (Glassmorphic) usando Tailwind CSS. Implementa animaciones fluidas con Framer Motion que reaccionan inmediatamente cuando el WebSocket emite fluctuaciones de latencia o caídas de red.

---

## GUÍA DE EJECUCIÓN PASO A PASO

### Preparación de la Base de Datos (Común para ambos OS)
Antes de arrancar los scripts, asegúrese de tener PostgreSQL corriendo y cree una base de datos limpia:
```sql
CREATE DATABASE pulseflux_db;