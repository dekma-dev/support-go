import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  AuthSession,
  Ticket,
  TicketPriority,
  createTicket,
  listTickets,
} from "./api";

export function TicketList({ session }: { session: AuthSession | null }) {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [priority, setPriority] = useState<TicketPriority>("medium");

  const loadTickets = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const items = await listTickets();
      setTickets(items);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load tickets");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadTickets();
  }, [loadTickets]);

  const canSubmit = useMemo(
    () => !submitting && !!session && title.trim().length > 0 && description.trim().length > 0,
    [description, session, submitting, title],
  );

  async function onCreateTicket(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!canSubmit) return;

    setSubmitting(true);
    setError(null);
    try {
      await createTicket({
        title: title.trim(),
        description: description.trim(),
        priority,
      });
      setTitle("");
      setDescription("");
      setPriority("medium");
      await loadTickets();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create ticket");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="card">
      <h2>Tickets</h2>

      {session ? (
        <form className="ticket-form" onSubmit={onCreateTicket}>
          <h3>Create ticket</h3>
          <label>
            Title
            <input
              value={title}
              onChange={(event) => setTitle(event.target.value)}
              placeholder="Payment page is unavailable"
              maxLength={200}
            />
          </label>

          <label>
            Description
            <textarea
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              placeholder="Steps to reproduce and observed behavior"
              rows={3}
              maxLength={4000}
            />
          </label>

          <div className="form-row">
            <label>
              Priority
              <select
                value={priority}
                onChange={(event) => setPriority(event.target.value as TicketPriority)}
              >
                <option value="low">low</option>
                <option value="medium">medium</option>
                <option value="high">high</option>
                <option value="urgent">urgent</option>
              </select>
            </label>
            <button type="submit" disabled={!canSubmit} style={{ alignSelf: "end" }}>
              {submitting ? "Creating..." : "Create ticket"}
            </button>
          </div>
        </form>
      ) : (
        <p className="muted">Login to create tickets.</p>
      )}

      <div className="tickets-head">
        <h3>All tickets</h3>
        <button type="button" onClick={() => void loadTickets()} disabled={loading}>
          {loading ? "Refreshing..." : "Refresh"}
        </button>
      </div>

      {error ? <p className="error">{error}</p> : null}
      {!loading && tickets.length === 0 ? <p className="muted">No tickets yet.</p> : null}

      <ul className="ticket-list">
        {tickets.map((ticket) => (
          <li key={ticket.id}>
            <Link to={`/tickets/${ticket.id}`} className="ticket-link">
              <article>
                <div className="ticket-meta">
                  <strong>{ticket.public_id}</strong>
                  <span className={`badge badge-status-${ticket.status}`}>{ticket.status}</span>
                  <span className={`badge badge-priority-${ticket.priority}`}>{ticket.priority}</span>
                  {ticket.assignee_id ? (
                    <span className="badge">agent: {ticket.assignee_id.slice(0, 12)}</span>
                  ) : null}
                </div>
                <h4>{ticket.title}</h4>
                <p className="ticket-desc">{ticket.description}</p>
                <small className="muted">
                  Created {new Date(ticket.created_at).toLocaleString()}
                </small>
              </article>
            </Link>
          </li>
        ))}
      </ul>
    </section>
  );
}
