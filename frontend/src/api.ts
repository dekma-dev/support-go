export type TicketStatus =
  | "new"
  | "open"
  | "pending_customer"
  | "pending_internal"
  | "resolved"
  | "closed";

export type TicketPriority = "low" | "medium" | "high" | "urgent";
export type UserRole = "client" | "agent" | "admin";

export type AuthSession = {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_at: string;
  refresh_expires_at: string;
  user_id: string;
  email: string;
  role: UserRole;
};

export type Ticket = {
  id: string;
  public_id: string;
  title: string;
  description: string;
  status: TicketStatus;
  priority: TicketPriority;
  requester_id: string;
  assignee_id?: string;
  sla_due_at?: string;
  created_at: string;
  updated_at: string;
  closed_at?: string;
};

export type CreateTicketInput = {
  title: string;
  description: string;
  requester_id: string;
  priority: TicketPriority;
  sla_due_at?: string;
};

type APIErrorResponse = {
  error?: string;
};

const defaultAPIBaseURL = "http://localhost:8080";
const sessionStorageKey = "support-go-session";

export const apiBaseURL = (
  import.meta.env.VITE_API_BASE_URL || defaultAPIBaseURL
).replace(/\/+$/, "");

async function request<T>(
  path: string,
  init?: RequestInit,
  options?: { withAuth?: boolean },
): Promise<T> {
  const session = readStoredSession();
  const response = await fetch(`${apiBaseURL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options?.withAuth !== false && session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : {}),
      ...(init?.headers || {}),
    },
    ...init,
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    try {
      const body = (await response.json()) as APIErrorResponse;
      if (body.error) {
        message = body.error;
      }
    } catch {
      // Keep generic message if server does not return JSON.
    }
    throw new Error(message);
  }

  return (await response.json()) as T;
}

export function listTickets() {
  return request<Ticket[]>("/api/v1/tickets");
}

export function createTicket(input: CreateTicketInput) {
  return request<Ticket>("/api/v1/tickets", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function readStoredSession(): AuthSession | null {
  if (typeof window === "undefined") {
    return null;
  }

  const rawValue = window.localStorage.getItem(sessionStorageKey);
  if (!rawValue) {
    return null;
  }

  try {
    return JSON.parse(rawValue) as AuthSession;
  } catch {
    window.localStorage.removeItem(sessionStorageKey);
    return null;
  }
}

export function clearStoredSession() {
  if (typeof window !== "undefined") {
    window.localStorage.removeItem(sessionStorageKey);
  }
}

function persistSession(session: AuthSession) {
  if (typeof window !== "undefined") {
    window.localStorage.setItem(sessionStorageKey, JSON.stringify(session));
  }
}

export async function register(input: {
  email: string;
  password: string;
  role: UserRole;
}) {
  const session = await request<AuthSession>(
    "/api/v1/auth/register",
    {
      method: "POST",
      body: JSON.stringify(input),
    },
    { withAuth: false },
  );
  persistSession(session);
  return session;
}

export async function login(input: { email: string; password: string }) {
  const session = await request<AuthSession>(
    "/api/v1/auth/login",
    {
      method: "POST",
      body: JSON.stringify(input),
    },
    { withAuth: false },
  );
  persistSession(session);
  return session;
}

export async function refreshSession(refreshToken?: string) {
  const currentSession = readStoredSession();
  const token = refreshToken || currentSession?.refresh_token;
  if (!token) {
    throw new Error("No refresh token available");
  }

  const session = await request<AuthSession>("/api/v1/auth/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token: token }),
  }, { withAuth: false });
  persistSession(session);
  return session;
}
