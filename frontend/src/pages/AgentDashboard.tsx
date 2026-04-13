import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Icon } from "../components/Icon";
import { listTickets, type AuthSession, type Ticket } from "../api";

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  return `${Math.floor(hrs / 24)}d ago`;
}

export function AgentDashboard({ session }: { session: AuthSession }) {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try {
      setLoading(true);
      setTickets(await listTickets());
    } catch {
      // silent
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  const myTickets = tickets.filter((t) => t.assignee_id === session.user_id && t.status !== "resolved" && t.status !== "closed");
  const resolved = tickets.filter((t) => t.status === "resolved" || t.status === "closed");
  const urgent = myTickets.filter((t) => t.priority === "urgent");

  const priorityBorder = (p: string) => {
    if (p === "urgent") return "border-l-4 border-sc-error";
    return "border-l-4 border-sc-primary/40";
  };

  return (
    <div className="p-8 space-y-10">
      {/* Header */}
      <header className="flex items-end justify-between">
        <div>
          <h1 className="font-headline font-bold text-5xl tracking-tighter text-sc-on-surface uppercase">OPERATIONS CONTROL</h1>
          <div className="flex items-center gap-3 mt-2">
            <span className="h-1 w-12 bg-sc-primary rounded-full" />
            <p className="text-sc-on-surface-variant font-label text-xs tracking-[0.2em] uppercase">Real-Time Telemetry</p>
          </div>
        </div>
        <div className="text-right hidden md:block">
          <p className="font-headline text-3xl font-bold text-sc-primary">{new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</p>
          <p className="font-label text-[10px] text-sc-on-surface-variant uppercase tracking-widest">System Uptime</p>
        </div>
      </header>

      {/* Stats */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="glass-panel p-6 rounded-lg relative overflow-hidden group">
          <div className="relative z-10">
            <div className="flex justify-between items-start mb-4">
              <span className="font-label text-[10px] font-bold tracking-widest text-sc-on-surface-variant uppercase">Tickets Resolved</span>
              <Icon name="check_circle" className="text-sc-primary text-lg" />
            </div>
            <div className="flex items-baseline gap-2">
              <span className="font-headline text-4xl font-bold text-sc-primary">{resolved.length}</span>
            </div>
            <div className="mt-6 h-12 w-full flex items-end gap-1">
              {[10, 20, 40, 60, 80].map((h, i) => (
                <div key={i} className="w-full bg-sc-primary rounded-t-sm transition-all duration-500" style={{ height: `${h}%`, opacity: 0.1 + i * 0.2 }} />
              ))}
            </div>
          </div>
        </div>

        <div className="glass-panel p-6 rounded-lg relative overflow-hidden group">
          <div className="relative z-10">
            <div className="flex justify-between items-start mb-4">
              <span className="font-label text-[10px] font-bold tracking-widest text-sc-on-surface-variant uppercase">My Active</span>
              <Icon name="assignment" className="text-sc-secondary text-lg" />
            </div>
            <div className="flex items-baseline gap-2">
              <span className="font-headline text-4xl font-bold text-sc-secondary">{myTickets.length}</span>
            </div>
            <div className="mt-6 h-12 w-full flex items-end gap-1">
              {[30, 50, 20, 60, 40].map((h, i) => (
                <div key={i} className="w-full bg-sc-secondary rounded-t-sm transition-all duration-500" style={{ height: `${h}%`, opacity: 0.1 + i * 0.2 }} />
              ))}
            </div>
          </div>
        </div>

        <div className="glass-panel p-6 rounded-lg relative overflow-hidden group">
          <div className="relative z-10">
            <div className="flex justify-between items-start mb-4">
              <span className="font-label text-[10px] font-bold tracking-widest text-sc-on-surface-variant uppercase">Current Load</span>
              <Icon name="bolt" className="text-sc-tertiary text-lg" />
            </div>
            <div className="flex items-baseline gap-2">
              <span className="font-headline text-4xl font-bold text-sc-tertiary">
                {myTickets.length === 0 ? "Idle" : myTickets.length < 5 ? "Low" : myTickets.length < 10 ? "Medium" : "Heavy"}
              </span>
              <span className="text-sc-tertiary-fixed text-xs font-bold font-label">
                {urgent.length > 0 ? `${urgent.length} URGENT` : "OPTIMAL"}
              </span>
            </div>
            <div className="mt-6 h-12 w-full flex items-end gap-1">
              {[20, 20, 40, 20, 30].map((h, i) => (
                <div key={i} className="w-full bg-sc-tertiary rounded-t-sm transition-all duration-500" style={{ height: `${h}%`, opacity: 0.1 + i * 0.2 }} />
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Main grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Left: My tickets */}
        <div className="lg:col-span-8 space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="font-headline text-xl font-bold text-sc-on-surface flex items-center gap-2">
              <Icon name="assignment" className="text-sc-primary" />
              MY ASSIGNED TICKETS
            </h3>
            <span className="font-label text-[10px] text-sc-on-surface-variant/60 tracking-widest uppercase">{myTickets.length} Total Active</span>
          </div>

          {loading ? (
            <div className="text-sc-on-surface-variant text-sm py-12 text-center">Loading...</div>
          ) : myTickets.length === 0 ? (
            <div className="glass-panel p-8 text-center text-sc-on-surface-variant">No tickets assigned to you.</div>
          ) : (
            <div className="space-y-4">
              {myTickets.slice(0, 10).map((t) => (
                <Link key={t.id} to={`/tickets/${t.public_id}`} className={`block glass-panel p-5 ${priorityBorder(t.priority)} relative overflow-hidden group cursor-pointer transition-all hover:bg-sc-surface-high`}>
                  <div className="flex justify-between items-start">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-3 mb-1">
                        <span className={`font-label text-[10px] font-bold tracking-widest uppercase ${t.priority === "urgent" ? "text-sc-error" : "text-sc-on-surface-variant"}`}>
                          {t.priority.toUpperCase()} #{t.public_id}
                        </span>
                        {t.priority === "urgent" && <span className="h-1.5 w-1.5 bg-sc-error rounded-full animate-pulse" />}
                      </div>
                      <h4 className="font-headline text-lg font-bold text-sc-on-surface group-hover:text-sc-primary transition-colors">{t.title}</h4>
                      <p className="text-sm text-sc-on-surface-variant mt-1 line-clamp-1">{t.description}</p>
                    </div>
                    <div className="text-right shrink-0 ml-4">
                      <p className="font-label text-[10px] text-sc-on-surface-variant uppercase tracking-tighter">Assigned {timeAgo(t.updated_at)}</p>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Right sidebar */}
        <aside className="lg:col-span-4 space-y-8">
          {/* Quick search */}
          <div className="space-y-4">
            <h3 className="font-headline text-lg font-bold text-sc-on-surface uppercase tracking-tight">Quick Search</h3>
            <div className="relative">
              <input
                type="text"
                className="w-full bg-sc-surface-highest border-none focus:ring-1 focus:ring-sc-primary rounded-sm text-xs py-3 px-10 placeholder:text-sc-on-surface-variant/40 font-label tracking-widest text-sc-on-surface"
                placeholder="TICKET ID, USER, OR IP..."
              />
              <Icon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-sc-on-surface-variant/60" />
            </div>
          </div>

          {/* Notifications */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-headline text-lg font-bold text-sc-on-surface uppercase tracking-tight">Notifications</h3>
              {urgent.length > 0 && (
                <span className="px-2 py-0.5 bg-sc-error/20 text-sc-error text-[10px] font-bold rounded-sm">{urgent.length} URGENT</span>
              )}
            </div>
            <div className="space-y-3">
              {urgent.slice(0, 3).map((t) => (
                <Link key={t.id} to={`/tickets/${t.public_id}`} className="block p-4 bg-sc-surface-low rounded border-l-2 border-sc-error/50">
                  <div className="flex gap-4 items-start">
                    <Icon name="warning" className="text-sc-error shrink-0" />
                    <div>
                      <p className="text-xs font-bold text-sc-on-surface">{t.title}</p>
                      <p className="text-[10px] text-sc-on-surface-variant mt-1">Priority: urgent — #{t.public_id}</p>
                    </div>
                  </div>
                </Link>
              ))}
              {urgent.length === 0 && (
                <div className="p-4 bg-sc-surface-low rounded border-l-2 border-sc-tertiary/30 flex gap-4 items-start">
                  <Icon name="check_circle" className="text-sc-tertiary shrink-0" />
                  <div>
                    <p className="text-xs font-bold text-sc-on-surface">All Clear</p>
                    <p className="text-[10px] text-sc-on-surface-variant mt-1">No urgent tickets assigned.</p>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* System log */}
          <div className="space-y-4">
            <h3 className="font-headline text-lg font-bold text-sc-on-surface uppercase tracking-tight">System Log</h3>
            <div className="bg-sc-surface-lowest p-4 rounded h-48 overflow-y-auto relative">
              <div className="absolute left-0 top-0 bottom-0 w-[2px] bg-gradient-to-b from-sc-primary/50 to-transparent" />
              <div className="space-y-3 font-body text-[10px] tracking-wider leading-relaxed">
                {tickets.slice(0, 6).map((t, i) => (
                  <div key={t.id} className="flex gap-2">
                    <span className="text-sc-primary opacity-70 shrink-0">[{new Date(t.updated_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}]</span>
                    <span className={i === 0 ? "text-sc-tertiary" : "text-sc-on-surface-variant"}>
                      {t.status === "resolved" ? "RESOLVED" : t.status === "open" ? "UPDATED" : "STATUS"}: #{t.public_id}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
}
