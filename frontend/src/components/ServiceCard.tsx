import React from 'react';
import { Microservice } from '../hooks/usePulseFlux';

interface Props {
  service: Microservice;
}

export const ServiceCard: React.FC<Props> = ({ service }) => {
  const isOnline = service.status === 'ONLINE';

  // Generación dinámica de polilíneas SVG para dibujar el flujo de latencia
  const points = service.latency_history
    .map((lat, index) => `${index * 40},${Math.max(10, 60 - lat / 5)}`)
    .join(' ');

  return (
    <div className={`glass-panel p-6 rounded-xl transition-all duration-500 transform hover:-translate-y-1 ${
      isOnline ? 'neon-border-online' : 'neon-border-offline'
    }`}>
      <div className="flex justify-between items-start mb-4">
        <div>
          <h3 className="text-lg font-bold text-white tracking-wide">{service.name}</h3>
          <p className="text-xs text-slate-400 font-mono overflow-hidden text-ellipsis whitespace-nowrap max-w-[200px]">
            {service.url}
          </p>
        </div>
        
        {/* Alertas Visuales con pulsación nativa CSS */}
        <span className={`flex h-3 w-3 relative`}>
          <span className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${
            isOnline ? 'bg-emerald-400' : 'bg-rose-400'
          }`}></span>
          <span className={`relative inline-flex rounded-full h-3 w-3 ${
            isOnline ? 'bg-emerald-500' : 'bg-rose-500'
          }`}></span>
        </span>
      </div>

      {/* Gráfica Vectorial Animada (Flujo de Latencia) */}
      <div className="h-16 w-full bg-[#0d0e15]/80 rounded-lg p-2 font-mono overflow-hidden mb-4 border border-white/[0.02]">
        <svg className="w-full h-full" viewBox="0 0 280 60" preserveAspectRatio="none">
          <polyline
            fill="none"
            stroke={isOnline ? "#10b981" : "#ef4444"}
            strokeWidth="2"
            points={points}
            className="transition-all duration-300"
          />
        </svg>
      </div>

      <div className="flex justify-between items-center text-xs font-mono text-slate-400">
        <span>Intervalo: {service.check_interval}s</span>
        <span className={`font-bold ${isOnline ? 'text-emerald-400' : 'text-rose-400'}`}>
          {isOnline ? `${service.latency_history[service.latency_history.length - 1] || 0} ms` : 'DISCONNECTED'}
        </span>
      </div>
    </div>
  );
};