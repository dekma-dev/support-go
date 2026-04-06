import { Link, useLocation } from "react-router-dom";
import { Icon } from "./Icon";
import type { AuthSession } from "../api";
import { clearStoredSession } from "../api";

type NavItem = { to: string; icon: string; label: string; dividerAfter?: boolean };

const navItems: NavItem[] = [
  { to: "/dashboard", icon: "dashboard", label: "My Tickets" },
  { to: "/urgent", icon: "priority_high", label: "Urgent" },
  { to: "/tickets", icon: "confirmation_number", label: "All Tickets", dividerAfter: true },
  { to: "/users", icon: "group", label: "User Management" },
  { to: "/admin", icon: "admin_panel_settings", label: "Admin Tools" },
  { to: "/logs", icon: "terminal", label: "System Logs" },
];

export function Sidebar({
  session,
  onSessionChange,
}: {
  session: AuthSession | null;
  onSessionChange: (s: AuthSession | null) => void;
}) {
  const { pathname } = useLocation();

  const handleLogout = () => {
    clearStoredSession();
    onSessionChange(null);
  };

  return (
    <aside className="hidden lg:flex flex-col h-full w-64 border-r border-sc-outline-variant/15 bg-sc-surface-low shrink-0">
      {/* Identity block */}
      <div className="px-6 py-6 mb-2 flex items-center gap-3">
        <div className="w-10 h-10 flex items-center justify-center bg-sc-primary-container/10 rounded-lg">
          <Icon name="terminal" className="text-sc-primary-container" />
        </div>
        <div>
          <h2 className="font-headline font-bold text-sc-primary-container tracking-tight text-sm">OPERATIONS</h2>
          <p className="text-[10px] uppercase tracking-[0.2em] text-sc-on-surface-variant/60">
            {session ? session.role.toUpperCase() : "GUEST"}
          </p>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 space-y-1">
        {navItems.map((item) => {
          const active = pathname === item.to || (item.to !== "/" && pathname.startsWith(item.to));
          return (
            <div key={item.to}>
              <Link
                to={item.to}
                className={`flex items-center gap-3 px-3 py-2.5 rounded font-label font-semibold text-xs uppercase tracking-widest transition-all duration-200 ${
                  active
                    ? "bg-sc-primary-container/10 text-sc-primary-container border-r-4 border-sc-primary-container shadow-[inset_-10px_0_15px_-10px_rgba(0,229,255,0.3)]"
                    : "text-sc-on-surface-variant opacity-70 hover:opacity-100 hover:bg-sc-surface-container"
                }`}
              >
                <Icon name={item.icon} className="text-sm" />
                {item.label}
              </Link>
              {item.dividerAfter && <div className="my-3 border-t border-sc-outline-variant/10" />}
            </div>
          );
        })}
      </nav>

      {/* Create ticket + footer */}
      <div className="px-4 mt-auto pb-6 space-y-4">
        <Link
          to="/tickets/new"
          className="block w-full py-3 bg-sc-primary-fixed text-sc-on-primary-fixed font-headline font-bold text-sm tracking-widest uppercase rounded-sm shadow-[0_0_15px_rgba(0,218,243,0.3)] active:scale-95 transition-transform text-center"
        >
          CREATE TICKET
        </Link>
        <div className="pt-4 border-t border-sc-outline-variant/15 space-y-1">
          <a href="#" className="flex items-center gap-3 px-3 py-2 text-sc-on-surface-variant opacity-70 hover:opacity-100 transition-all font-label text-xs uppercase tracking-widest">
            <Icon name="help" className="text-sm" />
            Help
          </a>
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2 text-sc-on-surface-variant opacity-70 hover:opacity-100 transition-all font-label text-xs uppercase tracking-widest w-full"
          >
            <Icon name="logout" className="text-sm" />
            Sign Out
          </button>
        </div>
      </div>
    </aside>
  );
}
