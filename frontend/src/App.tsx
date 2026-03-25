import { useEffect, useState } from "react";
import { Link, Route, Routes } from "react-router-dom";
import { AuthSession, readStoredSession } from "./api";
import { SessionCard } from "./SessionCard";
import { TicketList } from "./TicketList";
import { TicketDetail } from "./TicketDetail";

function HomePage() {
  return (
    <section className="card">
      <h2>Support-Go</h2>
      <p>
        Helpdesk system with Go API, PostgreSQL, Kafka events, and React UI.
      </p>
      <p>
        <Link to="/tickets">Open tickets</Link> to get started.
      </p>
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
        <Link to="/" className="logo">Support-Go</Link>
        <nav className="nav">
          <Link to="/tickets">Tickets</Link>
          {session ? (
            <span className="nav-user">{session.email} ({session.role})</span>
          ) : null}
        </nav>
      </header>

      <SessionCard session={session} onSessionChange={setSession} />

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/tickets" element={<TicketList session={session} />} />
        <Route path="/tickets/:id" element={<TicketDetail session={session} />} />
      </Routes>
    </main>
  );
}
