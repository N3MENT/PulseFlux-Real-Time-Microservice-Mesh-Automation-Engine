@echo off
:: PulseFlux - Windows Automation Deployment Script
title PulseFlux Runner

echo =============================================
echo    PULSEFLUX - CONFIGURADOR NATIVO WINDOWS   
echo =============================================

:: Configuración de variables de entorno predeterminadas
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_NAME=pulseflux_db
set PORT=8080

echo [+] Preparando entorno para la compilación del Core...
cd /d "%~dp0..\backend"

:: Crear directorio dummy de distribución si no existe
if not exist dist mkdir dist
if not exist dist\index.html echo ^<h1^>PulseFlux UI en mantenimiento^</h1^> > dist\index.html

echo [+] Compilando ejecutable nativo para entorno Windows (.exe)...
go build -o ..\pulseflux_windows_amd64.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo [X] Error crítico: No se pudo compilar el binario del servidor.
    pause
    exit /b %ERRORLEVEL%
)

echo [+] Lanzando el servicio local PulseFlux...
cd ..
pulseflux_windows_amd64.exe
pause