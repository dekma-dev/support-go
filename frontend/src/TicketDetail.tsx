import { FormEvent, useCallback, useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { Icon } from "./components/Icon";
import {
  AuthSession, Comment, Ticket, TicketEvent, TicketStatus,
  addComment, assignTicket, changeTicketStatus, getTicket, listComments, listEvents,
} from "./api";

const ALL_STATUSES: TicketStatus[] = ["new", "open", "pending_customer", "pending_internal", "resolved", "closed"];

function timeLabel(dateStr: string) {
  return new Date(dateStr).toLocaleString(undefined, { hour: "2-digit", minute: "2-digit", hour12: true });
}

function dateLabel(dateStr: string) {
  const d = new Date(dateStr);
  const today = new Date();
  if (d.toDateString() === today.toDateString()) return "Today";
  return d.toLocaleDateString();
}

export function TicketDetail({ session }: { session: AuthSession | null }) {
  const { id } = useParams<{ id: string }>();
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [events, setEvents] = useState<TicketEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [tab, setTab] = useState<"comments" | "internal">("comments");
  const isStaff = session?.role === "agent" || session?.role === "admin";

  const load = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    setError(null);
    try {
      const [t, c, e] = await Promise.all([getTicket(id), listComments(id), listEvents(id)]);
      setTicket(t);
      setComments(c);
      setEvents(e);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load ticket");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => { void load(); }, [load]);

  if (loading) return <div className="flex items-center justify-center h-full text-sc-on-surface-variant">Loading...</div>;
  if (error) return <div className="p-8"><div className="text-sc-error mb-4">{error}</div><Link to="/tickets" className="text-sc-primary">Back to tickets</Link></div>;
  if (!ticket) return <div className="p-8 text-sc-on-surface-variant">Ticket not found</div>;

  const publicComments = comments.filter((c) => !c.is_internal);
  const internalNotes = comments.filter((c) => c.is_internal);
  const priorityColor = ticket.priority === "urgent" ? "text-sc-error" : ticket.priority === "high" ? "text-sc-primary-container" : "text-sc-on-surface-variant";

  return (
    <div className="flex flex-col md:flex-row h-full overflow-hidden">
      {/* Left: Event Timeline */}
      <section className="w-full md:w-80 lg:w-96 bg-sc-surface-low border-r border-sc-outline-variant/15 flex flex-col overflow-y-auto shrink-0">
        <div className="p-6 border-b border-sc-outline-variant/15">
          <h2 className="font-headline text-lg font-semibold text-sc-primary uppercase tracking-wider">Event Timeline</h2>
          <p className="text-xs text-sc-on-surface-variant font-label mt-1">REAL-TIME SEQUENTIAL LOG</p>
        </div>
        <div className="p-6 relative flex-1">
          <div className="absolute left-9 top-10 bottom-10 w-px bg-gradient-to-b from-sc-primary/40 via-sc-secondary/20 to-transparent" />
          <div className="space-y-8 relative">
            {events.map((evt, i) => (
              <div key={evt.id} className="flex gap-4">
                <div className="relative z-10 flex items-center justify-center w-6 h-6 rounded-full bg-sc-surface-highest border border-sc-primary/50 text-[10px] text-sc-primary shrink-0">
                  {String(i + 1).padStart(2, "0")}
                </div>
                <div className="min-w-0">
                  <p className="text-xs font-bold text-sc-on-surface uppercase">{evt.event_type.replace(/[._]/g, " ")}</p>
                  <p className="text-[10px] text-sc-on-surface-variant uppercase mt-0.5">
                    {dateLabel(evt.created_at)}, {timeLabel(evt.created_at)}
                  </p>
                  {evt.new_value && (
                    <p className="text-xs text-sc-primary/80 mt-2 bg-sc-primary/5 p-2 rounded border border-sc-primary/10 break-all">
                      {JSON.stringify(evt.new_value)}
                    </p>
                  )}
                </div>
              </div>
            ))}
            {events.length === 0 && (
              <p className="text-xs text-sc-on-surface-variant">No events recorded yet.</p>
            )}
          </div>
        </div>
      </section>

      {/* Center: Ticket Info & Comments */}
      <section className="flex-1 flex flex-col bg-sc-surface overflow-y-auto min-w-0">
        {/* Header bar */}
        <div className="p-6 border-b border-sc-outline-variant/15 flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <Link to="/tickets" className="text-sc-on-surface-variant hover:text-sc-primary transition-colors">
              <Icon name="arrow_back" className="text-lg" />
            </Link>
            <h1 className="font-headline text-3xl font-bold text-sc-on-surface tracking-tighter">#{ticket.public_id}</h1>
            <div>
              <span className={`bg-sc-error/20 ${priorityColor} px-2 py-0.5 text-[10px] font-bold tracking-widest border border-current/30 uppercase`}>
                {ticket.priority}
              </span>
              <p className="text-sm font-headline text-sc-primary mt-1">{ticket.title}</p>
            </div>
          </div>
        </div>

        <div className="p-6 space-y-8">
          {/* Info cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-5 glass-card glow-border flex flex-col">
              <p className="text-[10px] text-sc-on-surface-variant font-bold uppercase tracking-widest">Requester</p>
              <div className="flex items-center gap-3 mt-4">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-sc-primary to-sc-secondary p-0.5">
                  <div className="w-full h-full rounded-full bg-sc-surface flex items-center justify-center font-bold text-sc-primary text-sm">
                    {ticket.requester_id.slice(0, 2).toUpperCase()}
                  </div>
                </div>
                <p className="text-xs font-mono text-sc-on-surface-variant">{ticket.requester_id.slice(0, 16)}</p>
              </div>
            </div>
            <div className="p-5 glass-card glow-border flex flex-col">
              <p className="text-[10px] text-sc-on-surface-variant font-bold uppercase tracking-widest">Assigned To</p>
              <p className="text-sm font-bold text-sc-on-surface mt-4">
                {ticket.assignee_id ? ticket.assignee_id.slice(0, 16) : <span className="italic text-sc-on-surface-variant/40">Unassigned</span>}
              </p>
            </div>
            <div className="p-5 glass-card glow-border flex flex-col">
              <p className="text-[10px] text-sc-on-surface-variant font-bold uppercase tracking-widest">Status</p>
              <p className="text-xl font-headline font-bold text-sc-primary mt-4 uppercase">{ticket.status.replace("_", " ")}</p>
            </div>
          </div>

          {/* Description */}
          <div className="space-y-4">
            <h3 className="text-xs font-bold text-sc-on-surface uppercase tracking-widest flex items-center gap-2">
              <span className="w-1.5 h-1.5 bg-sc-primary rounded-full" />
              Original Report
            </h3>
            <div className="p-6 bg-sc-surface-container/50 border border-sc-outline-variant/10 leading-relaxed text-sm text-sc-on-surface/90">
              {ticket.description}
            </div>
          </div>

          {/* Tabs */}
          <div className="space-y-6">
            <div className="flex items-center border-b border-sc-outline-variant/15 gap-8">
              <button
                className={`pb-3 text-xs font-bold uppercase tracking-widest ${tab === "comments" ? "text-sc-primary border-b-2 border-sc-primary" : "text-sc-on-surface-variant hover:text-sc-on-surface"} transition-colors`}
                onClick={() => setTab("comments")}
              >
                Public Comments ({publicComments.length})
              </button>
              {isStaff && (
                <button
                  className={`pb-3 text-xs font-bold uppercase tracking-widest ${tab === "internal" ? "text-sc-primary border-b-2 border-sc-primary" : "text-sc-on-surface-variant hover:text-sc-on-surface"} transition-colors`}
                  onClick={() => setTab("internal")}
                >
                  Internal Notes ({internalNotes.length})
                </button>
              )}
            </div>

            {/* Comment list */}
            <div className="space-y-4">
              {(tab === "comments" ? publicComments : internalNotes).map((c) => (
                <div key={c.id} className="flex gap-4">
                  <div className="w-10 h-10 rounded-full bg-sc-surface-highest border border-sc-primary/20 flex items-center justify-center text-xs font-bold text-sc-primary shrink-0">
                    {c.author_id.slice(0, 2).toUpperCase()}
                  </div>
                  <div className="flex-1">
                    <div className={`glass-card glow-border p-4 ${c.is_internal ? "bg-sc-secondary/5 border-sc-secondary/10" : ""}`}>
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-xs font-bold text-sc-on-surface">
                          {c.author_id.slice(0, 16)} <span className="text-[10px] font-normal text-sc-on-surface-variant ml-2">{timeLabel(c.created_at)}</span>
                        </span>
                        {c.is_internal && <span className="text-[10px] bg-sc-secondary/20 text-sc-secondary px-2 py-0.5 uppercase tracking-widest font-bold">Internal</span>}
                      </div>
                      <p className="text-sm text-sc-on-surface/80">{c.body}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Comment editor */}
            {session && (
              <CommentEditor ticketId={ticket.id} isStaff={isStaff} onAdded={load} />
            )}
          </div>
        </div>
      </section>

      {/* Right: State Controls */}
      {isStaff && (
        <section className="w-full md:w-64 lg:w-72 bg-sc-surface-low border-l border-sc-outline-variant/15 p-6 space-y-8 overflow-y-auto shrink-0">
          <StaffPanel ticket={ticket} session={session!} onUpdate={load} />
        </section>
      )}
    </div>
  );
}

function StaffPanel({ ticket, session, onUpdate }: { ticket: Ticket; session: AuthSession; onUpdate: () => void }) {
  const [status, setStatus] = useState<TicketStatus>(ticket.status);
  const [assigneeId, setAssigneeId] = useState(ticket.assignee_id || "");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => { setStatus(ticket.status); setAssigneeId(ticket.assignee_id || ""); }, [ticket]);

  async function doAssign() {
    if (!assigneeId.trim()) return;
    setBusy(true); setError(null);
    try { await assignTicket(ticket.id, assigneeId.trim()); onUpdate(); }
    catch (e) { setError(e instanceof Error ? e.message : "Failed"); }
    finally { setBusy(false); }
  }

  async function doStatus(s: TicketStatus) {
    setBusy(true); setError(null);
    try { await changeTicketStatus(ticket.id, s); onUpdate(); }
    catch (e) { setError(e instanceof Error ? e.message : "Failed"); }
    finally { setBusy(false); }
  }

  return (
    <>
      <div>
        <h3 className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest mb-4">State Transition</h3>
        <div className="space-y-3">
          <button onClick={() => doStatus("resolved")} disabled={busy} className="w-full py-3 bg-sc-tertiary/10 text-sc-tertiary border border-sc-tertiary/20 hover:bg-sc-tertiary hover:text-sc-on-tertiary transition-all font-bold text-xs flex items-center justify-center gap-2 uppercase disabled:opacity-50">
            <Icon name="check_circle" className="text-sm" /> Resolve Ticket
          </button>
          <button onClick={() => doStatus("closed")} disabled={busy} className="w-full py-3 bg-sc-surface-highest text-sc-on-surface border border-sc-outline-variant/20 hover:bg-sc-surface-bright transition-all font-bold text-xs flex items-center justify-center gap-2 uppercase disabled:opacity-50">
            <Icon name="lock" className="text-sm" /> Close Permanently
          </button>
        </div>
      </div>

      <div>
        <h3 className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest mb-4">Assign Agent</h3>
        <div className="space-y-2">
          <input
            value={assigneeId}
            onChange={(e) => setAssigneeId(e.target.value)}
            placeholder="User ID"
            className="w-full bg-sc-surface-highest border-none text-xs rounded-sm py-2 px-3 text-sc-on-surface placeholder:text-sc-on-surface-variant/40 focus:ring-1 focus:ring-sc-primary"
          />
          <div className="flex gap-2">
            <button onClick={() => { setAssigneeId(session.user_id); }} className="flex-1 py-2 bg-sc-primary/10 text-sc-primary border border-sc-primary/20 text-[10px] font-bold uppercase">Me</button>
            <button onClick={doAssign} disabled={busy || !assigneeId.trim()} className="flex-1 py-2 bg-sc-primary-fixed text-sc-on-primary-fixed text-[10px] font-bold uppercase disabled:opacity-50">Assign</button>
          </div>
        </div>
      </div>

      <div>
        <h3 className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest mb-4">Change Status</h3>
        <select
          value={status}
          onChange={(e) => { setStatus(e.target.value as TicketStatus); doStatus(e.target.value as TicketStatus); }}
          className="w-full bg-sc-surface-highest border-none text-xs rounded-sm py-2 px-3 text-sc-on-surface focus:ring-1 focus:ring-sc-primary"
        >
          {ALL_STATUSES.map((s) => <option key={s} value={s}>{s.replace("_", " ")}</option>)}
        </select>
      </div>

      {error && <div className="text-sc-error text-xs">{error}</div>}

      <div>
        <h3 className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest mb-4">Metadata</h3>
        <div className="space-y-3 text-xs">
          <div><p className="text-[10px] text-sc-on-surface-variant/60">CREATED</p><p className="font-mono text-sc-on-surface">{new Date(ticket.created_at).toLocaleString()}</p></div>
          <div><p className="text-[10px] text-sc-on-surface-variant/60">UPDATED</p><p className="font-mono text-sc-on-surface">{new Date(ticket.updated_at).toLocaleString()}</p></div>
          {ticket.closed_at && <div><p className="text-[10px] text-sc-on-surface-variant/60">CLOSED</p><p className="font-mono text-sc-on-surface">{new Date(ticket.closed_at).toLocaleString()}</p></div>}
        </div>
      </div>
    </>
  );
}

function CommentEditor({ ticketId, isStaff, onAdded }: { ticketId: string; isStaff: boolean; onAdded: () => void }) {
  const [body, setBody] = useState("");
  const [isInternal, setIsInternal] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!body.trim()) return;
    setSubmitting(true); setError(null);
    try {
      await addComment(ticketId, { body: body.trim(), is_internal: isInternal });
      setBody(""); setIsInternal(false); onAdded();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed");
    } finally { setSubmitting(false); }
  }

  return (
    <form onSubmit={onSubmit} className="p-4 bg-sc-surface-high/40 rounded-lg border border-sc-outline-variant/15">
      <textarea
        value={body}
        onChange={(e) => setBody(e.target.value)}
        placeholder="Type your response..."
        className="w-full bg-transparent border-none focus:ring-0 text-sm h-24 resize-none text-sc-on-surface placeholder:text-sc-on-surface-variant/40"
      />
      <div className="flex items-center justify-between mt-4">
        <div className="flex items-center gap-4">
          {isStaff && (
            <label className="flex items-center gap-2 text-xs text-sc-on-surface-variant cursor-pointer">
              <input type="checkbox" checked={isInternal} onChange={(e) => setIsInternal(e.target.checked)} className="rounded-sm bg-sc-surface-container border-sc-outline-variant text-sc-primary focus:ring-sc-primary/20" />
              Internal Note
            </label>
          )}
        </div>
        <button
          type="submit"
          disabled={submitting || !body.trim()}
          className="bg-sc-primary text-sc-on-primary px-6 py-2 text-xs font-bold uppercase tracking-widest hover:shadow-[0_0_20px_rgba(195,245,255,0.4)] transition-all disabled:opacity-50"
        >
          {submitting ? "Sending..." : "Submit Action"}
        </button>
      </div>
      {error && <div className="text-sc-error text-xs mt-2">{error}</div>}
    </form>
  );
}
