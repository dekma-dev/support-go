import { Link } from "react-router-dom";
import { Icon } from "../components/Icon";

export function DashboardHome() {
  return (
    <div className="relative">
      {/* Hero glow */}
      <div className="absolute inset-0 pointer-events-none" style={{ background: "radial-gradient(circle at 50% 50%, rgba(0,229,255,0.15) 0%, rgba(11,19,38,0) 70%)" }} />

      <section className="max-w-6xl mx-auto px-8 pt-20 pb-32 relative z-10">
        <div className="flex flex-col items-center text-center space-y-8">
          <div className="space-y-4">
            <h1 className="text-6xl md:text-7xl font-bold font-headline tracking-tighter text-sc-on-surface leading-none">
              Welcome to <span className="text-transparent bg-clip-text bg-gradient-to-r from-sc-primary-fixed to-sc-secondary">Command Center</span>
            </h1>
            <p className="text-xl text-sc-on-surface-variant/80 font-body max-w-2xl mx-auto">
              Your primary interface for operational assistance and technical support. Engineered for precision and speed.
            </p>
          </div>

          {/* Search bar */}
          <div className="w-full max-w-3xl glass-panel p-1 rounded-lg shadow-2xl glow-border">
            <div className="flex items-center bg-sc-surface-highest/40 rounded-sm px-6 py-4">
              <Icon name="search" className="text-sc-primary/60 mr-4" />
              <input
                type="text"
                className="bg-transparent border-none focus:ring-0 w-full text-sc-on-surface placeholder:text-sc-on-surface-variant/40 font-body text-lg"
                placeholder="Search system knowledge, tickets, or documentation..."
              />
              <div className="flex items-center gap-2">
                <span className="px-2 py-1 bg-sc-surface-container rounded text-[10px] text-sc-on-surface-variant font-mono border border-sc-outline-variant/20">CMD</span>
                <span className="px-2 py-1 bg-sc-surface-container rounded text-[10px] text-sc-on-surface-variant font-mono border border-sc-outline-variant/20">K</span>
              </div>
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex items-center gap-6 pt-4">
            <Link
              to="/tickets/new"
              className="px-8 py-4 bg-sc-primary-fixed text-sc-on-primary-fixed font-headline font-bold rounded-sm shadow-[0_0_20px_rgba(0,218,243,0.3)] hover:shadow-[0_0_30px_rgba(0,218,243,0.5)] transition-all active:scale-95 flex items-center gap-2"
            >
              <Icon name="add_task" />
              CREATE TICKET
            </Link>
            <Link
              to="/tickets"
              className="text-sc-on-surface font-headline font-semibold hover:text-sc-primary transition-colors flex items-center gap-2 group"
            >
              VIEW MY TICKETS
              <Icon name="arrow_forward_ios" className="text-sm group-hover:translate-x-1 transition-transform" />
            </Link>
          </div>
        </div>

        {/* Info cards */}
        <div className="mt-32 grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="glass-panel p-8 rounded-sm border border-sc-outline-variant/10 glow-border group hover:bg-sc-surface-low transition-colors">
            <div className="w-12 h-12 bg-sc-primary/10 rounded flex items-center justify-center mb-6 text-sc-primary">
              <Icon name="auto_graph" className="text-3xl" />
            </div>
            <h3 className="text-xl font-headline font-bold text-sc-on-surface mb-3">Live System Status</h3>
            <p className="text-sc-on-surface-variant text-sm leading-relaxed mb-6">Real-time monitoring of global infrastructure and node health.</p>
            <div className="flex items-center gap-2">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-sc-tertiary opacity-75" />
                <span className="relative inline-flex rounded-full h-2 w-2 bg-sc-tertiary" />
              </span>
              <span className="text-[10px] font-bold text-sc-tertiary uppercase tracking-widest">All Systems Operational</span>
            </div>
          </div>

          <div className="glass-panel p-8 rounded-sm border border-sc-outline-variant/10 glow-border group hover:bg-sc-surface-low transition-colors">
            <div className="w-12 h-12 bg-sc-secondary/10 rounded flex items-center justify-center mb-6 text-sc-secondary">
              <Icon name="history" className="text-3xl" />
            </div>
            <h3 className="text-xl font-headline font-bold text-sc-on-surface mb-3">Recent Activity</h3>
            <p className="text-sc-on-surface-variant text-sm leading-relaxed mb-6">Review your last interactions and status updates on open cases.</p>
            <div className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest">3 Active Threads</div>
          </div>

          <div className="glass-panel p-8 rounded-sm border border-sc-outline-variant/10 glow-border group hover:bg-sc-surface-low transition-colors">
            <div className="w-12 h-12 bg-sc-tertiary/10 rounded flex items-center justify-center mb-6 text-sc-tertiary">
              <Icon name="menu_book" className="text-3xl" />
            </div>
            <h3 className="text-xl font-headline font-bold text-sc-on-surface mb-3">Documentation</h3>
            <p className="text-sc-on-surface-variant text-sm leading-relaxed mb-6">Browse our extensive technical library for self-service resolution.</p>
            <div className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest">1.2k Protocols Available</div>
          </div>
        </div>
      </section>
    </div>
  );
}
