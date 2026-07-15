#!/bin/bash
# PulseFlux - Linux Automation Deployment Script
set -e

echo "============================================="
echo "   PULSEFLUX - CONFIGURADOR NATIVO LINUX    "
echo "============================================="

# Configuración por defecto de variables de entorno para PostgreSQL
export DB_HOST=${DB_HOST:-"localhost"}
export DB_PORT=${DB_PORT:-"5432"}
export DB_USER=${DB_USER:-"postgres"}
export DB_NAME=${DB_NAME:-"pulseflux_db"}
export PORT=${PORT:-"8080"}

# Verificar la disponibilidad local de PostgreSQL
echo "[+] Verificando conexión local con PostgreSQL..."
if ! command -v pg_isready &> /dev/null; then
    echo "[!] Advertencia: 'pg_isready' no está instalado. Continuando con la ejecución directa..."
else
    until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
      echo "[!] Esperando a que el servicio PostgreSQL se estabilice..."
      sleep 2
    done
    echo "[+] Conexión con PostgreSQL validada con éxito."
fi

# Compilar backend embebiendo la UI de producción
echo "[+] Compilando binario monolítico ejecutable de Linux..."
cd "$(dirname "$0")/../backend"

# Creamos una carpeta temporal "dist" simulada si no existe para que compile Go
mkdir -p dist
if [ ! -f dist/index.html ]; then
    echo "<h1>PulseFlux UI en mantenimiento</h1>" > dist/index.html
fi

go build -o ../pulseflux_linux_amd64 main.go

echo "[+] Iniciando aplicación web PulseFlux..."
cd ..
./pulseflux_linux_amd64