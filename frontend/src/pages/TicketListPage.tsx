import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Icon } from "../components/Icon";
import { listTickets, type AuthSession, type ListTicketsFilter, type Ticket, type TicketPriority, type TicketStatus } from "../api";

type FilterMode = "all" | "high_priority" | "mine";

function statusBadge(status: TicketStatus) {
  const map: Record<TicketStatus, { bg: string; text: string; glow: string }> = {
    new: { bg: "bg-sc-primary/10 border-sc-primary/20", text: "text-sc-primary", glow: "status-glow-open" },
    open: { bg: "bg-sc-error/10 border-sc-error/20", text: "text-sc-error", glow: "status-glow-urgent" },
    pending_customer: { bg: "bg-sc-secondary-container/20 border-sc-secondary-container/30", text: "text-sc-on-secondary-container", glow: "" },
    pending_internal: { bg: "bg-sc-secondary-container/20 border-sc-secondary-container/30", text: "text-sc-on-secondary-container", glow: "" },
    resolved: { bg: "bg-sc-tertiary-fixed/10 border-sc-tertiary-fixed/20", text: "text-sc-tertiary-fixed", glow: "status-glow-resolved" },
    closed: { bg: "bg-sc-surface-high border-sc-outline-variant/20", text: "text-sc-on-surface-variant", glow: "" },
  };
  const s = map[status] ?? map.new;
  return (
    <span className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full ${s.bg} ${s.text} ${s.glow} text-[10px] font-bold tracking-tighter uppercase border`}>
      <span className={`w-1 h-1 rounded-full ${s.text.replace("text-", "bg-")} ${status === "open" ? "animate-pulse" : ""}`} />
      {status.replace("_", " ")}
    </span>
  );
}

function priorityBadge(p: TicketPriority) {
  if (p === "urgent") {
    return (
      <div className="relative inline-block px-3 py-1 rounded-sm border-2 border-sc-error/30 bg-sc-error/5 overflow-hidden">
        <span className="text-[10px] font-bold text-sc-error uppercase tracking-widest relative z-10">Urgent</span>
        <div className="absolute inset-0 bg-sc-error/10 animate-pulse" />
      </div>
    );
  }
  if (p === "high") {
    return <span className="text-[10px] font-bold text-sc-primary-container uppercase tracking-widest px-2 py-1 bg-sc-primary-container/10 rounded-sm">High</span>;
  }
  return <span className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest px-2 py-1 bg-sc-surface-high rounded-sm">{p}</span>;
}

function timeAgo(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ${mins % 60}m ago`;
  return `${Math.floor(hrs / 24)}d ago`;
}

export function TicketListPage({ session }: { session: AuthSession | null }) {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [mode, setMode] = useState<FilterMode>("all");

  const load = useCallback(async () => {
    try {
      setLoading(true);
      setError("");

      const filter: ListTicketsFilter = { sort: "created_at_desc", limit: 100 };
      if (mode === "high_priority") {
        filter.priority = ["high", "urgent"];
      } else if (mode === "mine" && session) {
        filter.assignee_id = "me";
      }

      const response = await listTickets(filter);
      setTickets(response.items);
      setTotal(response.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load tickets");
    } finally {
      setLoading(false);
    }
  }, [mode, session]);

  useEffect(() => { load(); }, [load]);

  const openCount = tickets.filter((t) => t.status !== "resolved" && t.status !== "closed").length;

  const filterBtnClass = (active: boolean) =>
    `px-4 py-2 font-label font-bold text-[10px] tracking-widest uppercase rounded-sm border flex items-center gap-2 transition-colors ${
      active
        ? "bg-sc-surface-high text-sc-primary border-sc-primary/20"
        : "bg-sc-surface-lowest text-sc-on-surface-variant hover:text-sc-primary border-sc-outline-variant/15"
    }`;

  return (
    <div className="p-6 space-y-6">
      {/* Header + Stats */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <span className="w-2 h-2 rounded-full bg-sc-primary animate-pulse" />
            <span className="text-[10px] font-label font-bold tracking-[0.3em] text-sc-primary uppercase">Active Command Session</span>
          </div>
          <h2 className="text-4xl font-headline font-bold text-sc-on-surface tracking-tight">TICKET_REGISTRY</h2>
        </div>
        <div className="flex gap-4">
          <div className="glass-card shimmer-border px-6 py-3 rounded-lg flex flex-col">
            <span className="text-[10px] font-label font-bold text-sc-on-surface-variant/60 tracking-widest uppercase">Total Open</span>
            <span className="text-2xl font-headline font-bold text-sc-primary">{openCount}</span>
          </div>
          <div className="glass-card shimmer-border px-6 py-3 rounded-lg flex flex-col">
            <span className="text-[10px] font-label font-bold text-sc-on-surface-variant/60 tracking-widest uppercase">Total</span>
            <span className="text-2xl font-headline font-bold text-sc-tertiary">{total}</span>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3">
        <button onClick={() => setMode("all")} className={filterBtnClass(mode === "all")}>
          <Icon name="filter_list" className="text-xs" />
          All Tickets
        </button>
        <button onClick={() => setMode("high_priority")} className={filterBtnClass(mode === "high_priority")}>
          High Priority
        </button>
        {session && (session.role === "agent" || session.role === "admin") && (
          <button onClick={() => setMode("mine")} className={filterBtnClass(mode === "mine")}>
            Assigned to Me
          </button>
        )}
        <div className="h-4 w-px bg-sc-outline-variant/30 mx-2" />
        <span className="text-[10px] font-label font-bold text-sc-on-surface-variant/40 tracking-widest uppercase">Sorted by: Latest</span>
      </div>

      {/* Table */}
      {error && <div className="text-sc-error text-sm">{error}</div>}
      {loading ? (
        <div className="text-sc-on-surface-variant text-sm py-20 text-center">Loading tickets...</div>
      ) : (
        <div className="glass-card rounded-xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-sc-surface-high/50 border-b border-sc-outline-variant/15">
                  <th className="px-6 py-4 text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-[0.2em]">Ticket ID</th>
                  <th className="px-6 py-4 text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-[0.2em]">Subject</th>
                  <th className="px-6 py-4 text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-[0.2em]">Status</th>
                  <th className="px-6 py-4 text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-[0.2em]">Priority</th>
                  <th className="px-6 py-4 text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-[0.2em]">Created</th>
                  <th className="px-6 py-4" />
                </tr>
              </thead>
              <tbody className="divide-y divide-sc-outline-variant/10">
                {tickets.map((t) => (
                  <tr key={t.id} className={`group hover:bg-sc-surface-high/30 transition-colors ${t.status === "resolved" || t.status === "closed" ? "opacity-60" : ""}`}>
                    <td className="px-6 py-5 font-headline font-bold text-sc-primary-fixed tracking-widest text-xs">
                      {t.public_id}
                    </td>
                    <td className="px-6 py-5">
                      <Link to={`/tickets/${t.public_id}`} className="flex flex-col">
                        <span className={`text-sm font-semibold text-sc-on-surface group-hover:text-sc-primary transition-colors ${t.status === "resolved" || t.status === "closed" ? "line-through" : ""}`}>
                          {t.title}
                        </span>
                        <span className="text-[10px] text-sc-on-surface-variant/70 uppercase tracking-wider line-clamp-1">
                          {t.description}
                        </span>
                      </Link>
                    </td>
                    <td className="px-6 py-5">{statusBadge(t.status)}</td>
                    <td className="px-6 py-5">{priorityBadge(t.priority)}</td>
                    <td className="px-6 py-5 text-xs text-sc-on-surface-variant/80 font-label">{timeAgo(t.created_at)}</td>
                    <td className="px-6 py-5 text-right">
                      <Link to={`/tickets/${t.public_id}`} className="p-2 text-sc-on-surface-variant hover:text-sc-primary transition-colors">
                        <Icon name="open_in_new" className="text-sm" />
                      </Link>
                    </td>
                  </tr>
                ))}
                {tickets.length === 0 && (
                  <tr>
                    <td colSpan={6} className="px-6 py-12 text-center text-sc-on-surface-variant">
                      No tickets found. Create your first ticket to get started.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* FAB */}
      <Link
        to="/tickets/new"
        className="fixed bottom-8 right-8 w-16 h-16 rounded-sm bg-sc-primary-fixed text-sc-on-primary-fixed shadow-[0_12px_24px_-8px_rgba(0,218,243,0.5)] hover:shadow-[0_0_25px_rgba(0,218,243,0.6)] flex items-center justify-center active:scale-90 transition-all z-50 group"
      >
        <Icon name="add" className="text-3xl group-hover:rotate-90 transition-transform duration-300" />
      </Link>
    </div>
  );
}
