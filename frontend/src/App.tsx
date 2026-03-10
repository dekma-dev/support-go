import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { Link, Route, Routes } from "react-router-dom";
import {
  apiBaseURL,
  createTicket,
  listTickets,
  Ticket,
  TicketPriority,
} from "./api";

function HomePage() {
  return (
    <section className="card">
      <h2>Support-Go UI</h2>
      <p>Frontend baseline is connected to the ticket API.</p>
    </section>
  );
}

function TicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [requesterID, setRequesterID] = useState("");
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
    () =>
      !submitting &&
      title.trim().length > 0 &&
      description.trim().length > 0 &&
      requesterID.trim().length > 0,
    [description, requesterID, submitting, title],
  );

  async function onCreateTicket(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!canSubmit) {
      return;
    }

    setSubmitting(true);
    setError(null);
    try {
      await createTicket({
        title: title.trim(),
        description: description.trim(),
        requester_id: requesterID.trim(),
        priority,
      });
      setTitle("");
      setDescription("");
      setRequesterID("");
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
      <p>API: {apiBaseURL}</p>

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
            rows={4}
            maxLength={4000}
          />
        </label>

        <div className="form-row">
          <label>
            Requester ID
            <input
              value={requesterID}
              onChange={(event) => setRequesterID(event.target.value)}
              placeholder="client_123"
              maxLength={100}
            />
          </label>

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
        </div>

        <button type="submit" disabled={!canSubmit}>
          {submitting ? "Creating..." : "Create ticket"}
        </button>
      </form>

      <div className="tickets-head">
        <h3>Ticket list</h3>
        <button type="button" onClick={() => void loadTickets()} disabled={loading}>
          {loading ? "Refreshing..." : "Refresh"}
        </button>
      </div>

      {error ? <p className="error">{error}</p> : null}
      {!loading && tickets.length === 0 ? <p>No tickets yet.</p> : null}

      <ul className="ticket-list">
        {tickets.map((ticket) => (
          <li key={ticket.id}>
            <article>
              <div className="ticket-meta">
                <strong>{ticket.public_id}</strong>
                <span className="badge">{ticket.status}</span>
                <span className="badge badge-priority">{ticket.priority}</span>
              </div>
              <h4>{ticket.title}</h4>
              <p>{ticket.description}</p>
              <small>Requester: {ticket.requester_id}</small>
            </article>
          </li>
        ))}
      </ul>
    </section>
  );
}

export function App() {
  return (
    <main className="layout">
      <header className="header">
        <h1>Support-Go</h1>
        <nav className="nav">
          <Link to="/">Home</Link>
          <Link to="/tickets">Tickets</Link>
        </nav>
      </header>

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/tickets" element={<TicketsPage />} />
      </Routes>
    </main>
  );
}
