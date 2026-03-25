import { FormEvent, useCallback, useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import {
  AuthSession,
  Comment,
  Ticket,
  TicketEvent,
  TicketStatus,
  addComment,
  assignTicket,
  changeTicketStatus,
  getTicket,
  listComments,
  listEvents,
} from "./api";

const ALL_STATUSES: TicketStatus[] = [
  "new",
  "open",
  "pending_customer",
  "pending_internal",
  "resolved",
  "closed",
];

export function TicketDetail({ session }: { session: AuthSession | null }) {
  const { id } = useParams<{ id: string }>();
  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [events, setEvents] = useState<TicketEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [tab, setTab] = useState<"comments" | "events">("comments");

  const isStaff = session?.role === "agent" || session?.role === "admin";

  const load = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    setError(null);
    try {
      const [t, c, e] = await Promise.all([
        getTicket(id),
        listComments(id),
        listEvents(id),
      ]);
      setTicket(t);
      setComments(c);
      setEvents(e);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load ticket");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    void load();
  }, [load]);

  if (loading) return <section className="card"><p>Loading...</p></section>;
  if (error) return <section className="card"><p className="error">{error}</p><Link to="/tickets">Back to list</Link></section>;
  if (!ticket) return <section className="card"><p>Ticket not found</p></section>;

  return (
    <>
      <section className="card">
        <Link to="/tickets" className="back-link">Back to tickets</Link>
        <div className="ticket-detail-header">
          <div>
            <h2>{ticket.public_id}: {ticket.title}</h2>
            <div className="ticket-meta">
              <span className={`badge badge-status-${ticket.status}`}>{ticket.status}</span>
              <span className={`badge badge-priority-${ticket.priority}`}>{ticket.priority}</span>
              {ticket.assignee_id ? (
                <span className="badge">assigned: {ticket.assignee_id.slice(0, 16)}</span>
              ) : (
                <span className="badge badge-empty">unassigned</span>
              )}
            </div>
          </div>
        </div>

        <p className="ticket-description">{ticket.description}</p>

        <div className="detail-meta-grid">
          <div><strong>Requester:</strong> {ticket.requester_id.slice(0, 16)}</div>
          <div><strong>Created:</strong> {new Date(ticket.created_at).toLocaleString()}</div>
          <div><strong>Updated:</strong> {new Date(ticket.updated_at).toLocaleString()}</div>
          {ticket.closed_at ? (
            <div><strong>Closed:</strong> {new Date(ticket.closed_at).toLocaleString()}</div>
          ) : null}
        </div>

        {isStaff ? (
          <StaffActions ticket={ticket} session={session!} onUpdate={load} />
        ) : null}
      </section>

      <section className="card">
        <div className="tab-bar">
          <button
            className={tab === "comments" ? "tab-active" : "tab-inactive"}
            onClick={() => setTab("comments")}
          >
            Comments ({comments.length})
          </button>
          <button
            className={tab === "events" ? "tab-active" : "tab-inactive"}
            onClick={() => setTab("events")}
          >
            History ({events.length})
          </button>
        </div>

        {tab === "comments" ? (
          <CommentsSection
            ticketId={ticket.id}
            comments={comments}
            session={session}
            onAdded={load}
          />
        ) : (
          <EventsSection events={events} />
        )}
      </section>
    </>
  );
}

function StaffActions({
  ticket,
  session,
  onUpdate,
}: {
  ticket: Ticket;
  session: AuthSession;
  onUpdate: () => void;
}) {
  const [assigneeId, setAssigneeId] = useState(ticket.assignee_id || "");
  const [status, setStatus] = useState<TicketStatus>(ticket.status);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setAssigneeId(ticket.assignee_id || "");
    setStatus(ticket.status);
  }, [ticket]);

  async function onAssign() {
    if (!assigneeId.trim()) return;
    setBusy(true);
    setError(null);
    try {
      await assignTicket(ticket.id, assigneeId.trim());
      onUpdate();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Assign failed");
    } finally {
      setBusy(false);
    }
  }

  async function onChangeStatus() {
    if (status === ticket.status) return;
    setBusy(true);
    setError(null);
    try {
      await changeTicketStatus(ticket.id, status);
      onUpdate();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Status change failed");
    } finally {
      setBusy(false);
    }
  }

  function selfAssign() {
    setAssigneeId(session.user_id);
  }

  return (
    <div className="staff-actions">
      <h4>Actions</h4>
      {error ? <p className="error">{error}</p> : null}
      <div className="form-row">
        <label>
          Assignee
          <div className="input-with-btn">
            <input
              value={assigneeId}
              onChange={(e) => setAssigneeId(e.target.value)}
              placeholder="user_id of agent"
            />
            <button type="button" className="button-small" onClick={selfAssign} disabled={busy}>
              Me
            </button>
          </div>
        </label>
        <button type="button" onClick={() => void onAssign()} disabled={busy || !assigneeId.trim()}>
          Assign
        </button>
      </div>
      <div className="form-row">
        <label>
          Status
          <select value={status} onChange={(e) => setStatus(e.target.value as TicketStatus)}>
            {ALL_STATUSES.map((s) => (
              <option key={s} value={s}>{s}</option>
            ))}
          </select>
        </label>
        <button
          type="button"
          onClick={() => void onChangeStatus()}
          disabled={busy || status === ticket.status}
        >
          Change status
        </button>
      </div>
    </div>
  );
}

function CommentsSection({
  ticketId,
  comments,
  session,
  onAdded,
}: {
  ticketId: string;
  comments: Comment[];
  session: AuthSession | null;
  onAdded: () => void;
}) {
  const [body, setBody] = useState("");
  const [isInternal, setIsInternal] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const isStaff = session?.role === "agent" || session?.role === "admin";

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!body.trim()) return;
    setSubmitting(true);
    setError(null);
    try {
      await addComment(ticketId, { body: body.trim(), is_internal: isInternal });
      setBody("");
      setIsInternal(false);
      onAdded();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add comment");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="comments-section">
      {comments.length === 0 ? (
        <p className="muted">No comments yet.</p>
      ) : (
        <ul className="comment-list">
          {comments.map((c) => (
            <li key={c.id} className={c.is_internal ? "comment-internal" : ""}>
              <div className="comment-header">
                <strong>{c.author_id.slice(0, 16)}</strong>
                <span className="muted">{new Date(c.created_at).toLocaleString()}</span>
                {c.is_internal ? <span className="badge badge-internal">internal</span> : null}
              </div>
              <p>{c.body}</p>
            </li>
          ))}
        </ul>
      )}

      {session ? (
        <form className="comment-form" onSubmit={onSubmit}>
          <textarea
            value={body}
            onChange={(e) => setBody(e.target.value)}
            placeholder="Write a comment..."
            rows={3}
            maxLength={4000}
          />
          <div className="comment-form-actions">
            {isStaff ? (
              <label className="checkbox-label">
                <input
                  type="checkbox"
                  checked={isInternal}
                  onChange={(e) => setIsInternal(e.target.checked)}
                />
                Internal note
              </label>
            ) : null}
            <button type="submit" disabled={submitting || !body.trim()}>
              {submitting ? "Sending..." : "Add comment"}
            </button>
          </div>
          {error ? <p className="error">{error}</p> : null}
        </form>
      ) : (
        <p className="muted">Login to comment.</p>
      )}
    </div>
  );
}

function EventsSection({ events }: { events: TicketEvent[] }) {
  if (events.length === 0) {
    return <p className="muted">No history yet.</p>;
  }

  return (
    <ul className="event-list">
      {events.map((e) => (
        <li key={e.id}>
          <div className="event-header">
            <span className="badge">{e.event_type}</span>
            <span className="muted">{new Date(e.created_at).toLocaleString()}</span>
            <span className="muted">by {e.actor_id.slice(0, 16)}</span>
          </div>
          {e.old_value || e.new_value ? (
            <div className="event-diff">
              {e.old_value ? <span className="diff-old">{JSON.stringify(e.old_value)}</span> : null}
              {e.new_value ? <span className="diff-new">{JSON.stringify(e.new_value)}</span> : null}
            </div>
          ) : null}
        </li>
      ))}
    </ul>
  );
}
