import React, { useEffect, useState, useRef, useCallback } from 'react';
import ForceGraph2D from 'react-force-graph-2d';
import { Activity, GitBranch, Zap, Server, Radio, Layers, Clock, Shield } from 'lucide-react';

// ------------------------------------------------------------------
// AESTHETIC CONFIGURATION
// ------------------------------------------------------------------
const NODE_WIDTH = 160;
const NODE_HEIGHT = 60;
const COLORS = {
  bg: '#02040a',
  grid: 'rgba(59, 130, 246, 0.05)',
  healthy: { fill: '#0f172a', border: '#3b82f6', text: '#bfdbfe', shadow: '#3b82f6' },
  warning: { fill: '#2a1205', border: '#f97316', text: '#fed7aa', shadow: '#f97316' },
  critical: { fill: '#2f0808', border: '#ef4444', text: '#fecaca', shadow: '#ef4444' }
};

const App = () => {
  const [graphData, setGraphData] = useState({ nodes: [], links: [] });
  const [events, setEvents] = useState([]);
  const [healthMap, setHealthMap] = useState({});
  
  // Responsive Graph Sizing
  const containerRef = useRef(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });
  const fgRef = useRef();

  // ------------------------------------------------------------------
  // 1. RESPONSIVE OBSERVER
  // ------------------------------------------------------------------
  useEffect(() => {
    // Updates graph size when window/container resizes
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { width, height } = entry.contentRect;
        setDimensions({ width, height });
      }
    });

    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }

    return () => resizeObserver.disconnect();
  }, []);

  // ------------------------------------------------------------------
  // 2. DATA HANDLING
  // ------------------------------------------------------------------
  useEffect(() => {
    const eventSource = new EventSource('http://localhost:8080/service');
    eventSource.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);
        if (payload.type === 'reverse_dependency') {
          buildGraphFromReverseMap(payload.data);
        } else {
          handleFailureEvent(payload);
        }
      } catch (err) {
        console.error("Error parsing JSON", err);
      }
    };
    return () => eventSource.close();
  }, []);

  const buildGraphFromReverseMap = (reverseMap) => {
    const nodes = new Set();
    const links = [];
    Object.entries(reverseMap).forEach(([provider, consumers]) => {
      nodes.add(provider);
      if (consumers) {
        consumers.forEach(consumer => {
          nodes.add(consumer);
          links.push({ source: consumer, target: provider });
        });
      }
    });
    setGraphData({
      nodes: Array.from(nodes).map(id => ({ id, group: 'core' })),
      links: links
    });
  };

  const handleFailureEvent = (payload) => {
    setEvents((prev) => [payload, ...prev].slice(0, 50));
    
    const newHealth = {};
    const root = payload.Root || payload.root;
    const failedNodes = payload.FailedNodes || payload.failednodes || [];

    if (root) newHealth[root] = 'critical';
    failedNodes.forEach(node => newHealth[node] = 'warning');

    setHealthMap(newHealth);
    setTimeout(() => setHealthMap({}), 4000);
  };

  // ------------------------------------------------------------------
  // 3. CANVAS RENDERING
  // ------------------------------------------------------------------
  const paintNode = useCallback((node, ctx, globalScale) => {
    const status = healthMap[node.id] || 'healthy';
    const theme = COLORS[status];
    
    const x = node.x - NODE_WIDTH / 2;
    const y = node.y - NODE_HEIGHT / 2;
    const radius = 8;

    // Glow
    if (status !== 'healthy') {
      ctx.shadowBlur = 30;
      ctx.shadowColor = theme.shadow;
    } else {
      ctx.shadowBlur = 0;
    }

    // Card
    ctx.beginPath();
    ctx.fillStyle = theme.fill;
    ctx.strokeStyle = theme.border;
    ctx.lineWidth = status === 'critical' ? 2 : 1;
    ctx.roundRect(x, y, NODE_WIDTH, NODE_HEIGHT, radius);
    ctx.fill();
    ctx.shadowBlur = 0;
    ctx.stroke();

    // Accent Line
    ctx.beginPath();
    ctx.fillStyle = theme.border;
    ctx.roundRect(x, y, 4, NODE_HEIGHT, { topLeft: radius, bottomLeft: radius });
    ctx.fill();

    // Status Dot
    ctx.beginPath();
    ctx.fillStyle = theme.border;
    ctx.arc(x + NODE_WIDTH - 12, y + 12, 3, 0, 2 * Math.PI);
    ctx.fill();

    // Text
    if (globalScale > 0.5) {
      ctx.font = `bold 14px "Segoe UI", Roboto, Helvetica, Arial, sans-serif`;
      ctx.textAlign = 'left';
      ctx.textBaseline = 'middle';
      ctx.fillStyle = theme.text;
      ctx.fillText(node.id, x + 16, y + NODE_HEIGHT / 3);
      
      ctx.font = `10px "Segoe UI", sans-serif`;
      ctx.fillStyle = theme.text;
      ctx.globalAlpha = 0.6;
      const statusText = status === 'critical' ? 'CRITICAL FAILURE' : status === 'warning' ? 'DEGRADED' : 'OPERATIONAL';
      ctx.fillText(statusText, x + 16, y + (NODE_HEIGHT / 3) * 2);
      ctx.globalAlpha = 1.0;
    }
  }, [healthMap]);

  return (
    // MAIN CONTAINER: Stack vertical on mobile, row on desktop (lg)
    <div className="flex flex-col lg:flex-row h-screen w-full bg-[#02040a] text-white font-sans overflow-hidden relative selection:bg-blue-500/30">
      
      {/* BACKGROUND PATTERN */}
      <div 
        className="absolute inset-0 z-0 pointer-events-none opacity-20"
        style={{
          backgroundImage: `
            linear-gradient(to right, ${COLORS.grid} 1px, transparent 1px),
            linear-gradient(to bottom, ${COLORS.grid} 1px, transparent 1px)
          `,
          backgroundSize: '40px 40px'
        }}
      />

      {/* LEFT SECTION: GRAPH 
        - Mobile: Height 55%
        - Desktop: Flex-1 (Takes remaining space)
      */}
      <div ref={containerRef} className="relative z-0 h-[55vh] lg:h-auto lg:flex-1 w-full overflow-hidden">
        
        {/* Floating Toolbar (Responsive sizing) */}
        <div className="absolute top-4 left-4 z-10 flex flex-col gap-2 lg:gap-4 pointer-events-none">
            <div className="bg-slate-900/80 backdrop-blur-xl p-3 lg:p-5 rounded-2xl border border-white/10 shadow-2xl flex items-center gap-3 lg:gap-4 animate-in fade-in slide-in-from-top-4 duration-700">
                <div className="p-2 lg:p-3 bg-blue-500/10 rounded-xl border border-blue-500/20 shadow-[0_0_15px_rgba(59,130,246,0.2)]">
                    <GitBranch className="w-5 h-5 lg:w-6 lg:h-6 text-blue-400" />
                </div>
                <div>
                    <h1 className="font-bold text-sm lg:text-lg text-white tracking-wide">Service Mesh</h1>
                    <div className="flex items-center gap-2 text-[10px] lg:text-xs text-slate-400 font-mono">
                      <span className="w-1.5 h-1.5 lg:w-2 lg:h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                      LIVE TOPOLOGY
                    </div>
                </div>
            </div>

            <div className="bg-slate-900/80 backdrop-blur-md p-2 lg:p-3 rounded-xl border border-white/10 shadow-xl flex gap-3 lg:gap-6 text-[10px] lg:text-xs font-medium text-slate-300">
                <div className="flex items-center gap-2"><div className="w-2 h-2 bg-blue-500 rounded-full shadow-[0_0_8px_#3b82f6]"></div> <span className="hidden sm:inline">Healthy</span></div>
                <div className="flex items-center gap-2"><div className="w-2 h-2 bg-orange-500 rounded-full shadow-[0_0_8px_#f97316]"></div> <span className="hidden sm:inline">Impacted</span></div>
                <div className="flex items-center gap-2"><div className="w-2 h-2 bg-red-500 rounded-full shadow-[0_0_8px_#ef4444]"></div> <span className="hidden sm:inline">Critical</span></div>
            </div>
        </div>

        <ForceGraph2D
          ref={fgRef}
          width={dimensions.width}
          height={dimensions.height}
          graphData={graphData}
          dagMode="lr"
          dagLevelDistance={220}
          backgroundColor="rgba(0,0,0,0)"
          nodeCanvasObject={paintNode}
          nodePointerAreaPaint={(node, color, ctx) => {
            ctx.fillStyle = color;
            ctx.fillRect(node.x - NODE_WIDTH / 2, node.y - NODE_HEIGHT / 2, NODE_WIDTH, NODE_HEIGHT);
          }}
          linkColor={() => '#1e293b'}
          linkWidth={2}
          linkDirectionalParticles={4}
          linkDirectionalParticleSpeed={0.005}
          linkDirectionalParticleWidth={2}
          linkDirectionalParticleColor={() => '#60a5fa'}
          d3AlphaDecay={0.05} 
          d3VelocityDecay={0.4}
          cooldownTicks={100}
        />
      </div>

      {/* RIGHT SECTION: LOGS 
        - Mobile: Height 45% (w-full)
        - Desktop: Width 420px (h-full)
      */}
      <div className="h-[45vh] lg:h-auto w-full lg:w-[420px] relative z-20 flex flex-col bg-slate-950/50 backdrop-blur-xl border-t lg:border-t-0 lg:border-l border-white/5 shadow-2xl">
        
        {/* Header */}
        <div className="p-3 lg:p-5 border-b border-white/5 bg-slate-900/50 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <Activity className="w-4 h-4 lg:w-5 lg:h-5 text-indigo-400" /> 
            <span className="font-bold text-sm lg:text-base text-slate-200 tracking-wide">Event Stream</span>
          </div>
          <div className="px-2 py-1 bg-indigo-500/10 border border-indigo-500/20 rounded text-[10px] text-indigo-300 font-mono flex items-center gap-1">
             <Radio className="w-3 h-3 animate-pulse" /> LISTENING
          </div>
        </div>

        {/* Feed */}
        <div className="flex-1 overflow-y-auto p-3 lg:p-4 space-y-3 lg:space-y-4 custom-scrollbar">
            {events.length === 0 && (
                <div className="h-full flex flex-col items-center justify-center text-slate-600 gap-4 opacity-50">
                    <div className="p-3 lg:p-4 bg-slate-900 rounded-full border border-slate-800">
                      <Server className="w-6 h-6 lg:w-8 lg:h-8" />
                    </div>
                    <p className="text-xs lg:text-sm font-mono">Awaiting system telemetry...</p>
                </div>
            )}
            
            {events.map((evt, idx) => {
                const root = evt.Root || evt.root;
                const failed = evt.FailedNodes || evt.failednodes || [];
                const radius = evt.BlastRadius || evt.blast_radius || 0;
                const time = evt.Time || evt.failed_time || new Date().toISOString();

                return (
                    <div 
                      key={idx} 
                      className="group relative bg-[#0b1221] hover:bg-[#131d33] border border-slate-800/50 hover:border-indigo-500/30 p-3 lg:p-4 rounded-xl transition-all duration-300 animate-in slide-in-from-right-8 fade-in shadow-lg"
                    >
                        {/* Connecting Line (Only on Desktop) */}
                        <div className="absolute -left-[21px] top-6 w-4 h-[1px] bg-slate-800 hidden xl:block"></div>
                        
                        {/* Header Row */}
                        <div className="flex justify-between items-start mb-2 lg:mb-3">
                            <div className="flex items-center gap-2">
                                <div className="p-1.5 bg-red-500/10 rounded-lg border border-red-500/20">
                                  <Zap className="w-3.5 h-3.5 lg:w-4 lg:h-4 text-red-400" />
                                </div>
                                <div>
                                  <p className="text-xs lg:text-sm font-bold text-slate-200">{root}</p>
                                  <p className="text-[9px] lg:text-[10px] text-red-400 font-medium">CRITICAL FAILURE</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-1 text-[10px] text-slate-500 font-mono bg-slate-900 px-2 py-1 rounded border border-slate-800">
                                <Clock className="w-3 h-3" />
                                {new Date(time).toLocaleTimeString()}
                            </div>
                        </div>
                        
                        {/* Impact Section */}
                        {failed.length > 0 ? (
                            <div className="bg-[#02040a] rounded-lg p-2 lg:p-3 border border-white/5">
                                <div className="flex justify-between items-center mb-2">
                                    <div className="flex items-center gap-1.5 text-[10px] lg:text-xs text-slate-400">
                                      <Layers className="w-3 h-3 lg:w-3.5 lg:h-3.5" />
                                      <span>Blast Radius</span>
                                    </div>
                                    <span className="text-[10px] lg:text-xs font-mono font-bold text-orange-400">
                                        {radius} Nodes
                                    </span>
                                </div>
                                <div className="flex flex-wrap gap-1 lg:gap-2">
                                    {failed.map((n, i) => (
                                        <span key={i} className="px-1.5 py-0.5 lg:px-2 lg:py-1 bg-orange-500/5 text-orange-300 border border-orange-500/10 rounded-md text-[10px] lg:text-[11px] font-medium">
                                            {n}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        ) : (
                           <div className="text-[10px] lg:text-xs text-emerald-500/70 bg-emerald-500/5 border border-emerald-500/10 p-2 rounded flex items-center justify-center gap-2">
                             <Shield className="w-3 h-3" /> Failure Isolated
                           </div>
                        )}
                    </div>
                );
            })}
        </div>
        
        {/* Footer Status Bar */}
        <div className="p-2 lg:p-3 bg-slate-900/80 border-t border-white/5 text-[9px] lg:text-[10px] text-slate-500 flex justify-between font-mono backdrop-blur-md">
           <span>SYSTEM: ONLINE</span>
           <span>LATENCY: 12ms</span>
        </div>
      </div>

    </div>
  );
};

export default App;