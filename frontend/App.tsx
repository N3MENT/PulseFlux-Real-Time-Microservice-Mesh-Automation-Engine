import React, { useState } from 'react';
import { usePulseFlux } from './hooks/usePulseFlux';
import { ServiceCard } from './components/ServiceCard';

export const App: React.FC = () => {
  const { services, loading, refetch } = usePulseFlux();
  const [name, setName] = useState('');
  const [url, setUrl] = useState('');
  const [interval, setInterval] = useState(5);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !url) return;

    try {
      const response = await fetch('/api/services', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, url, check_interval: Number(interval) }),
      });

      if (response.ok) {
        setName('');
        setUrl('');
        refetch(); // Actualizar malla de monitorización inmediatamente
      }
    } catch (err) {
      console.error("Fallo al registrar microservicio:", err);
    }
  };

  return (
    <div className="min-h-screen p-8 max-w-7xl mx-auto">
      {/* Encabezado Principal */}
      <header className="mb-12 border-b border-white/5 pb-6">
        <h1 className="text-3xl font-extrabold tracking-wider bg-gradient-to-r from-emerald-400 via-cyan-400 to-blue-500 bg-clip-text text-transparent">
          PULSEFLUX // CORE MESH
        </h1>
        <p className="text-slate-400 text-sm font-mono mt-1">REAL-TIME MONITORING & WORKFLOW PARALLEL ENGINE</p>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        {/* Panel Izquierdo: Formulario de Registro */}
        <div className="lg:col-span-1">
          <div className="glass-panel p-6 rounded-xl border border-cyan-500/10">
            <h2 className="text-md font-bold text-cyan-400 font-mono mb-4">// REGISTRAR_SERVICIO</h2>
            <form onSubmit={handleRegister} className="space-y-4">
              <div>
                <label className="block text-xs font-mono text-slate-400 mb-1">ALIAS INTERNO</label>
                <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Auth API" className="w-full bg-[#0d0e15] border border-white/10 rounded-md p-2 text-sm text-white focus:outline-none focus:border-cyan-500" />
              </div>
              <div>
                <label className="block text-xs font-mono text-slate-400 mb-1">ENDPOINT URL</label>
                <input type="url" value={url} onChange={e => setUrl(e.target.value)} placeholder="https://api.example.com/health" className="w-full bg-[#0d0e15] border border-white/10 rounded-md p-2 text-sm text-white focus:outline-none focus:border-cyan-500" />
              </div>
              <div>
                <label className="block text-xs font-mono text-slate-400 mb-1">SONDEO (SEGUNDOS)</label>
                <input type="number" value={interval} onChange={e => setInterval(Number(e.target.value))} min="2" className="w-full bg-[#0d0e15] border border-white/10 rounded-md p-2 text-sm text-white focus:outline-none focus:border-cyan-500" />
              </div>
              <button type="submit" className="w-full bg-cyan-600 hover:bg-cyan-500 text-white font-mono text-xs font-bold py-2.5 px-4 rounded-md transition-colors cursor-pointer shadow-lg shadow-cyan-600/10">
                AÑADIR A LA MALLA
              </button>
            </form>
          </div>
        </div>

        {/* Panel Derecho: Cuadrícula de Monitorización Reactiva */}
        <div className="lg:col-span-3">
          {loading ? (
            <div className="text-center font-mono text-sm text-slate-400 py-12 animate-pulse">Sincronizando canales mesh...</div>
          ) : services.length === 0 ? (
            <div className="glass-panel text-center text-slate-500 font-mono text-xs py-16 rounded-xl">
              NIGÚN ENLACE ACTIVO EN EL CLÚSTER. REGISTRE UN SERVICIO PARA INICIAR EL WORKER POOL.
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
              {services.map((svc) => (
                <ServiceCard key={svc.id} service={svc} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};