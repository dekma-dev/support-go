import { FormEvent, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Icon } from "../components/Icon";
import { createTicket, type TicketPriority } from "../api";

const PRIORITIES: { value: TicketPriority; label: string; color: string; desc: string }[] = [
  { value: "low", label: "Low", color: "text-sc-on-surface-variant border-sc-outline-variant/30", desc: "Non-blocking, routine" },
  { value: "medium", label: "Medium", color: "text-sc-primary border-sc-primary/30", desc: "Standard priority" },
  { value: "high", label: "High", color: "text-sc-primary-container border-sc-primary-container/30", desc: "Requires prompt attention" },
  { value: "urgent", label: "Urgent", color: "text-sc-error border-sc-error/30", desc: "Critical — immediate action" },
];

export function CreateTicketPage() {
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [priority, setPriority] = useState<TicketPriority>("medium");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!title.trim() || !description.trim()) return;
    setSubmitting(true);
    setError("");
    try {
      const ticket = await createTicket({ title: title.trim(), description: description.trim(), priority });
      navigate(`/tickets/${ticket.public_id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create ticket");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="p-8 max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-10">
        <Link to="/tickets" className="text-sc-on-surface-variant hover:text-sc-primary transition-colors text-xs font-label tracking-widest uppercase flex items-center gap-2 mb-4">
          <Icon name="arrow_back" className="text-sm" />
          Back to tickets
        </Link>
        <div className="flex items-center gap-2 mb-1">
          <span className="w-2 h-2 rounded-full bg-sc-primary animate-pulse" />
          <span className="text-[10px] font-label font-bold tracking-[0.3em] text-sc-primary uppercase">New Incident Report</span>
        </div>
        <h1 className="text-4xl font-headline font-bold text-sc-on-surface tracking-tight">CREATE_TICKET</h1>
      </div>

      <form onSubmit={onSubmit} className="space-y-8">
        {/* Title */}
        <div className="glass-card glow-border rounded-lg p-6 space-y-3">
          <label className="text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-widest flex items-center gap-2">
            <span className="w-1.5 h-1.5 bg-sc-primary rounded-full" />
            Subject
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Brief summary of the incident..."
            maxLength={200}
            required
            className="w-full bg-sc-surface-highest/50 border-none rounded-sm py-3 px-4 text-sc-on-surface placeholder:text-sc-on-surface-variant/40 focus:ring-1 focus:ring-sc-primary outline-none text-base"
          />
          <div className="flex justify-end">
            <span className="text-[10px] text-sc-on-surface-variant/50 font-label">{title.length} / 200</span>
          </div>
        </div>

        {/* Description */}
        <div className="glass-card glow-border rounded-lg p-6 space-y-3">
          <label className="text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-widest flex items-center gap-2">
            <span className="w-1.5 h-1.5 bg-sc-primary rounded-full" />
            Incident Description
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Detailed report: what happened, when, what you expected, relevant logs or error codes..."
            maxLength={4000}
            rows={10}
            required
            className="w-full bg-sc-surface-highest/50 border-none rounded-sm py-3 px-4 text-sc-on-surface placeholder:text-sc-on-surface-variant/40 focus:ring-1 focus:ring-sc-primary outline-none text-sm font-mono resize-none"
          />
          <div className="flex justify-end">
            <span className="text-[10px] text-sc-on-surface-variant/50 font-label">{description.length} / 4000</span>
          </div>
        </div>

        {/* Priority selector */}
        <div className="glass-card glow-border rounded-lg p-6 space-y-4">
          <label className="text-[10px] font-label font-bold text-sc-on-surface-variant uppercase tracking-widest flex items-center gap-2">
            <span className="w-1.5 h-1.5 bg-sc-primary rounded-full" />
            Priority Level
          </label>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            {PRIORITIES.map((p) => {
              const active = priority === p.value;
              return (
                <button
                  key={p.value}
                  type="button"
                  onClick={() => setPriority(p.value)}
                  className={`p-4 rounded-sm border-2 transition-all text-left ${
                    active
                      ? `${p.color} bg-current/5 shadow-[0_0_15px_rgba(0,229,255,0.15)]`
                      : "border-sc-outline-variant/20 text-sc-on-surface-variant hover:border-sc-outline-variant/50"
                  }`}
                >
                  <div className={`font-headline font-bold text-sm uppercase tracking-widest ${active ? "" : "text-sc-on-surface-variant"}`}>
                    {p.label}
                  </div>
                  <div className="text-[10px] mt-1 font-label opacity-70">{p.desc}</div>
                </button>
              );
            })}
          </div>
        </div>

        {/* Error */}
        {error && (
          <div className="text-sc-error text-sm font-label bg-sc-error/10 border border-sc-error/20 px-4 py-3 rounded-sm flex items-center gap-2">
            <Icon name="error" className="text-sm" />
            {error}
          </div>
        )}

        {/* Actions */}
        <div className="flex items-center justify-end gap-4 pt-4 border-t border-sc-outline-variant/15">
          <Link
            to="/tickets"
            className="px-6 py-3 bg-sc-surface-highest text-sc-on-surface border border-sc-outline-variant/20 hover:bg-sc-surface-bright transition-all font-headline font-bold text-xs uppercase tracking-widest"
          >
            Cancel
          </Link>
          <button
            type="submit"
            disabled={submitting || !title.trim() || !description.trim()}
            className="px-8 py-3 bg-sc-primary-fixed text-sc-on-primary-fixed font-headline font-bold text-xs uppercase tracking-widest shadow-[0_0_15px_rgba(0,218,243,0.3)] hover:shadow-[0_0_25px_rgba(0,218,243,0.5)] active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            {submitting ? (
              "SUBMITTING..."
            ) : (
              <>
                <Icon name="rocket_launch" className="text-sm" />
                INITIATE INCIDENT
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
}
