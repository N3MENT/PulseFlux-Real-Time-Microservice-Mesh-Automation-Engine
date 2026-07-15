-- Extensión para IDs únicos de forma nativa si es necesario
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tabla de Microservicios Monitoreados
CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL DEFAULT 'HTTP', -- HTTP, TCP, etc.
    check_interval INT NOT NULL DEFAULT 5,    -- Intervalo en segundos
    status VARCHAR(20) NOT NULL DEFAULT 'UNKNOWN', -- ONLINE, OFFLINE, UNKNOWN
    metadata JSONB DEFAULT '{}'::jsonb,       -- Cabeceras personalizadas o payloads
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabla de Historial y Métricas (Particionable en producción)
CREATE TABLE IF NOT EXISTS performance_logs (
    id BIGSERIAL PRIMARY KEY,
    service_id UUID REFERENCES services(id) ON DELETE CASCADE,
    latency_ms INT NOT NULL,
    status_code INT,
    error_message TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Índices para optimización de consultas de alta velocidad (Dashboard)
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON performance_logs(timestamp DESC);

-- 4. Trigger automático para actualizar el campo 'updated_at'
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_services_modtime
    BEFORE UPDATE ON services
    FOR EACH ROW
    EXECUTE PROCEDURE update_modified_column();