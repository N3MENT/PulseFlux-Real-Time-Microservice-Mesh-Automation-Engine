import { useState, useEffect } from 'react';

export interface Microservice {
  id: string;
  name: string;
  url: string;
  check_interval: number;
  status: 'ONLINE' | 'OFFLINE' | 'UNKNOWN';
  latency_history: number[];
}

export const usePulseFlux = () => {
  const [services, setServices] = useState<Microservice[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // 1. Cargar estado inicial mediante la API REST
  const fetchServices = async () => {
    try {
      const response = await fetch('/api/services');
      const data = await response.json();
      const initializedData = data.map((svc: any) => ({
        ...svc,
        latency_history: svc.status === 'ONLINE' ? [20] : [0] // Inicializador estático de gráficas
      }));
      setServices(initializedData || []);
    } catch (error) {
      console.error("Error consultando la API de control:", error);
    } finally {
      setLoading(false);
    }
  };

  // 2. Acoplar canal de transmisión síncrona en tiempo real (WebSockets)
  useEffect(() => {
    fetchServices();

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUri = `${protocol}//${window.location.host}/ws`;
    let socket = new WebSocket(wsUri);

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);
        if (payload.event === 'METRIC_UPDATE') {
          setServices((prevServices) =>
            prevServices.map((svc) => {
              if (svc.id === payload.service_id) {
                // Mantener los últimos 8 registros de latencia para las gráficas SVG animadas
                const updatedHistory = [...svc.latency_history, payload.latency_ms].slice(-8);
                return {
                  ...svc,
                  status: payload.status,
                  latency_history: updatedHistory,
                };
              }
              return svc;
            })
          );
        }
      } catch (err) {
        console.error("Error procesando trama WebSocket:", err);
      }
    };

    socket.onclose = () => {
      console.warn("Túnel WebSocket cerrado. Intentando reconexión pasiva...");
    };

    return () => socket.close();
  }, []);

  return { services, loading, refetch: fetchServices };
};