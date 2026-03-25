import { FormEvent, useEffect, useState } from "react";
import {
  AuthSession,
  UserRole,
  clearStoredSession,
  login,
  refreshSession,
  register,
} from "./api";

export function SessionCard({
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
        <span className={`badge ${session ? "badge-ok" : ""}`}>
          {session ? "Authenticated" : "Anonymous"}
        </span>
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
        </div>
      ) : null}
    </section>
  );
}
