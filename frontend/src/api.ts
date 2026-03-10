export type TicketStatus =
  | "new"
  | "open"
  | "pending_customer"
  | "pending_internal"
  | "resolved"
  | "closed";

export type TicketPriority = "low" | "medium" | "high" | "urgent";

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

export const apiBaseURL = (
  import.meta.env.VITE_API_BASE_URL || defaultAPIBaseURL
).replace(/\/+$/, "");

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBaseURL}${path}`, {
    headers: {
      "Content-Type": "application/json",
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
