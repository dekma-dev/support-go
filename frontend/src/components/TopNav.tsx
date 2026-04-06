import { Link, useLocation } from "react-router-dom";
import { Icon } from "./Icon";
import type { AuthSession } from "../api";

const topLinks = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/knowledge", label: "Knowledge Base" },
  { to: "/templates", label: "Templates" },
] as const;

export function TopNav({ session }: { session: AuthSession | null }) {
  const { pathname } = useLocation();

  return (
    <header className="bg-sc-bg border-b border-sc-outline-variant/15 shadow-[0_0_15px_rgba(0,229,255,0.1)] flex justify-between items-center w-full px-6 py-3 sticky top-0 z-50">
      <div className="flex items-center gap-8">
        <Link to="/" className="text-xl font-bold tracking-widest text-sc-primary-container uppercase font-headline">
          SYNTHETIC COMMAND
        </Link>
        <nav className="hidden md:flex items-center gap-6">
          {topLinks.map((l) => {
            const active = pathname.startsWith(l.to);
            return (
              <Link
                key={l.to}
                to={l.to}
                className={`text-sm uppercase tracking-widest font-semibold transition-colors ${
                  active
                    ? "text-sc-primary-container border-b-2 border-sc-primary-container pb-1"
                    : "text-sc-on-surface-variant hover:text-sc-primary-container"
                }`}
              >
                {l.label}
              </Link>
            );
          })}
        </nav>
      </div>

      <div className="flex items-center gap-4">
        <div className="relative hidden sm:block">
          <Icon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-sc-on-surface-variant text-sm" />
          <input
            type="text"
            placeholder="QUERY SYSTEM..."
            className="bg-sc-surface-highest/50 border-none rounded-sm pl-10 pr-4 py-1.5 text-xs font-label tracking-widest focus:ring-1 focus:ring-sc-primary w-64 transition-all text-sc-on-surface placeholder:text-sc-on-surface-variant/40"
          />
        </div>
        <div className="flex items-center gap-2">
          <button className="p-2 text-sc-on-surface-variant hover:bg-sc-surface-container rounded transition-all">
            <Icon name="notifications" />
          </button>
          <button className="p-2 text-sc-on-surface-variant hover:bg-sc-surface-container rounded transition-all">
            <Icon name="settings" />
          </button>
          {session && (
            <div className="h-8 w-8 rounded-full bg-sc-surface-high border border-sc-primary-container/20 flex items-center justify-center text-xs font-bold text-sc-primary-container ml-2">
              {session.email[0].toUpperCase()}
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
