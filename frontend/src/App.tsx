import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { Link, Route, Routes } from "react-router-dom";
import {
  AuthSession,
  UserRole,
  apiBaseURL,
  clearStoredSession,
  createTicket,
  listTickets,
  login,
  readStoredSession,
  refreshSession,
  register,
  Ticket,
  TicketPriority,
} from "./api";

function HomePage() {
  return (
    <section className="card">
      <h2>Support-Go UI</h2>
      <p>Frontend is connected to ticket API and real user auth backed by Postgres.</p>
    </section>
  );
}

function SessionCard({
  session,
  onSessionChange,
}: {
  session: AuthSession | null;
  onSessionChange: (session: AuthSession | null) => void;
}) {
  const [email, setEmail] = useState(session?.email || "");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<UserRole>(session?.role || "client");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setEmail(session?.email || "");
    setRole(session?.role || "client");
  }, [session]);

  async function onRegister(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const nextSession = await register({ email: email.trim(), password, role });
      setPassword("");
      onSessionChange(nextSession);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to register");
    } finally {
      setSubmitting(false);
    }
  }

  async function onLogin() {
    setSubmitting(true);
    setError(null);
    try {
      const nextSession = await login({ email: email.trim(), password });
      setPassword("");
      onSessionChange(nextSession);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to login");
    } finally {
      setSubmitting(false);
    }
  }

  async function onRefresh() {
    setSubmitting(true);
    setError(null);
    try {
      const nextSession = await refreshSession(session?.refresh_token);
      onSessionChange(nextSession);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to refresh session");
    } finally {
      setSubmitting(false);
    }
  }

  function onLogout() {
    clearStoredSession();
    setPassword("");
    onSessionChange(null);
  }

  return (
    <section className="card">
      <div className="session-head">
        <div>
          <h2>Session</h2>
          <p>Register or login with email/password, then use refresh rotation.</p>
        </div>
        <span className="badge">{session ? "Authenticated" : "Anonymous"}</span>
      </div>

      <form className="ticket-form" onSubmit={onRegister}>
        <div className="form-row">
          <label>
            Email
            <input
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="agent@example.com"
              maxLength={160}
            />
          </label>

          <label>
            Role
            <select value={role} onChange={(event) => setRole(event.target.value as UserRole)}>
              <option value="client">client</option>
              <option value="agent">agent</option>
              <option value="admin">admin</option>
            </select>
          </label>
        </div>

        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="Minimum 8 characters"
            maxLength={200}
          />
        </label>

        <div className="action-row">
          <button type="submit" disabled={submitting}>
            {submitting ? "Working..." : "Register"}
          </button>
          <button type="button" onClick={() => void onLogin()} disabled={submitting}>
            Login
          </button>
          <button type="button" onClick={() => void onRefresh()} disabled={submitting || !session}>
            Refresh
          </button>
          <button
            type="button"
            className="button-secondary"
            onClick={onLogout}
            disabled={submitting || !session}
          >
            Logout
          </button>
        </div>
      </form>

      {error ? <p className="error">{error}</p> : null}
      {session ? (
        <div className="session-grid">
          <p>
            <strong>User ID:</strong> {session.user_id}
          </p>
          <p>
            <strong>Email:</strong> {session.email}
          </p>
          <p>
            <strong>Role:</strong> {session.role}
          </p>
          <p>
            <strong>Access exp:</strong> {session.expires_at}
          </p>
          <p>
            <strong>Refresh exp:</strong> {session.refresh_expires_at}
          </p>
        </div>
      ) : null}
    </section>
  );
}

function TicketsPage({ session }: { session: AuthSession | null }) {
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
      <p>Auth: {session ? `${session.email} (${session.role})` : "no bearer token"}</p>

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
  const [session, setSession] = useState<AuthSession | null>(() => readStoredSession());

  useEffect(() => {
    setSession(readStoredSession());
  }, []);

  return (
    <main className="layout">
      <header className="header">
        <h1>Support-Go</h1>
        <nav className="nav">
          <Link to="/">Home</Link>
          <Link to="/tickets">Tickets</Link>
        </nav>
      </header>

      <SessionCard session={session} onSessionChange={setSession} />

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/tickets" element={<TicketsPage session={session} />} />
      </Routes>
    </main>
  );
}
