import { Icon } from "../components/Icon";
import type { AuthSession } from "../api";

// Note: backend doesn't have a list-users endpoint yet.
// This page is a UI shell that will be connected when the endpoint is added.

export function UserManagement({ session }: { session: AuthSession }) {
  const isAdmin = session.role === "admin";

  if (!isAdmin) {
    return (
      <div className="p-8">
        <div className="glass-panel p-12 text-center rounded-sm">
          <Icon name="lock" className="text-4xl text-sc-on-surface-variant/40 mb-4" />
          <h2 className="font-headline text-2xl font-bold text-sc-on-surface mb-2">ACCESS RESTRICTED</h2>
          <p className="text-sc-on-surface-variant text-sm">Admin clearance required for personnel management.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 space-y-10">
      {/* Header */}
      <header className="flex justify-between items-end">
        <div className="max-w-2xl">
          <h1 className="font-headline text-5xl font-bold tracking-tighter text-sc-on-surface mb-2">PERSONNEL MANAGEMENT</h1>
          <p className="text-sc-on-surface-variant text-sm font-light leading-relaxed">
            System administrator view for agent oversight, workload balancing, and tactical role assignment within the synthetic infrastructure.
          </p>
        </div>
        <button className="bg-sc-primary-container text-sc-on-primary-container px-6 py-3 rounded-sm font-bold tracking-widest text-xs flex items-center gap-3 shadow-[0_0_20px_rgba(0,229,255,0.2)] hover:shadow-[0_0_30px_rgba(0,229,255,0.4)] transition-all active:scale-95">
          <Icon name="person_add" />
          ADD NEW OPERATIVE
        </button>
      </header>

      {/* Stats */}
      <div className="grid grid-cols-12 gap-6">
        <div className="col-span-12 lg:col-span-4 glass-panel p-6 rounded-sm flex flex-col justify-between">
          <div>
            <div className="text-[10px] font-bold text-sc-primary tracking-[0.2em] mb-4 uppercase">Global Efficiency</div>
            <div className="font-headline text-6xl font-bold text-sc-on-surface">94.8<span className="text-sc-primary text-2xl">%</span></div>
          </div>
          <div className="h-1 bg-sc-surface-highest rounded-full overflow-hidden mt-4">
            <div className="h-full bg-sc-primary-container" style={{ width: "94.8%" }} />
          </div>
        </div>
        <div className="col-span-6 lg:col-span-4 glass-panel p-6 rounded-sm relative overflow-hidden group">
          <div className="relative z-10">
            <div className="text-[10px] font-bold text-sc-secondary tracking-[0.2em] mb-4 uppercase">Active Operatives</div>
            <div className="font-headline text-6xl font-bold text-sc-on-surface">3</div>
          </div>
          <Icon name="groups" className="absolute -bottom-4 -right-4 text-[96px] opacity-5 text-sc-secondary" />
        </div>
        <div className="col-span-6 lg:col-span-4 glass-panel p-6 rounded-sm border-l-2 border-sc-tertiary">
          <div className="text-[10px] font-bold text-sc-tertiary tracking-[0.2em] mb-4 uppercase">System Load</div>
          <div className="font-headline text-6xl font-bold text-sc-on-surface">Low</div>
          <div className="mt-4 flex gap-1">
            <div className="w-full h-1 bg-sc-tertiary/40" />
            <div className="w-full h-1 bg-sc-tertiary/40" />
            <div className="w-full h-1 bg-sc-surface-highest" />
            <div className="w-full h-1 bg-sc-surface-highest" />
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="glass-panel rounded-sm overflow-hidden relative">
        <div className="absolute top-0 left-0 w-full h-px" style={{ background: "linear-gradient(90deg, transparent, #00e5ff, transparent)" }} />
        <div className="px-6 py-4 border-b border-sc-outline-variant/10 flex justify-between items-center bg-sc-surface-low/50">
          <h3 className="font-headline text-lg font-medium text-sc-on-surface tracking-tight">Active Personnel Index</h3>
          <div className="flex items-center gap-2 px-3 py-1.5 bg-sc-surface-highest rounded-sm">
            <Icon name="filter_list" className="text-sm text-sc-on-surface-variant" />
            <span className="text-[10px] font-bold text-sc-on-surface-variant">SORT BY: RECENT</span>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-sc-surface-low/30">
                <th className="px-6 py-4 text-[10px] font-bold tracking-[0.1em] text-sc-on-surface-variant uppercase">Operative</th>
                <th className="px-6 py-4 text-[10px] font-bold tracking-[0.1em] text-sc-on-surface-variant uppercase">Role Authority</th>
                <th className="px-6 py-4 text-[10px] font-bold tracking-[0.1em] text-sc-on-surface-variant uppercase">Status</th>
                <th className="px-6 py-4 text-[10px] font-bold tracking-[0.1em] text-sc-on-surface-variant uppercase text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-sc-outline-variant/5">
              {/* Seed users - static for now until list-users endpoint exists */}
              <UserRow email="admin@synth.cmd" role="admin" roleColor="text-sc-primary" roleIcon="shield_person" />
              <UserRow email="agent@synth.cmd" role="agent" roleColor="text-sc-tertiary" roleIcon="star" />
              <UserRow email="client@synth.cmd" role="client" roleColor="text-sc-on-surface-variant" roleIcon="person" />
            </tbody>
          </table>
        </div>

        <div className="px-6 py-4 border-t border-sc-outline-variant/10 flex justify-between items-center bg-sc-surface-low/20">
          <div className="text-[10px] font-bold text-sc-on-surface-variant uppercase tracking-widest">Showing 3 operatives</div>
        </div>
      </div>

      {/* Footer */}
      <div className="flex justify-between items-center opacity-50 border-t border-sc-outline-variant/10 pt-6">
        <div className="text-[10px] font-medium tracking-[0.2em] text-sc-on-surface-variant flex items-center gap-4">
          <span>SECURITY STATUS: ENCRYPTED</span>
          <span>PROTOCOL: OPS-99</span>
        </div>
      </div>
    </div>
  );
}

function UserRow({ email, role, roleColor, roleIcon }: { email: string; role: string; roleColor: string; roleIcon: string }) {
  return (
    <tr className="hover:bg-sc-primary-container/5 transition-colors group">
      <td className="px-6 py-5">
        <div className="flex items-center gap-4">
          <div className="relative">
            <div className="w-10 h-10 rounded-sm bg-sc-surface-highest flex items-center justify-center text-sm font-bold text-sc-primary grayscale group-hover:grayscale-0 transition-all">
              {email[0].toUpperCase()}
            </div>
            <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-sc-primary-container rounded-full border-2 border-sc-surface" />
          </div>
          <div>
            <div className="text-sm font-semibold text-sc-on-surface">{email.split("@")[0]}</div>
            <div className="text-[10px] text-sc-on-surface-variant">{email}</div>
          </div>
        </div>
      </td>
      <td className="px-6 py-5">
        <div className="flex items-center gap-2">
          <Icon name={roleIcon} className={`${roleColor} text-lg`} />
          <span className={`text-xs font-bold ${roleColor} bg-current/10 px-2 py-0.5 rounded-sm uppercase`} style={{ backgroundColor: "currentColor", WebkitBackgroundClip: "unset", color: "inherit" }}>
            {/* Fix: use explicit bg */}
          </span>
          <span className={`text-xs font-bold uppercase ${roleColor}`}>{role}</span>
        </div>
      </td>
      <td className="px-6 py-5">
        <span className="text-[10px] text-sc-on-surface-variant uppercase">Active</span>
      </td>
      <td className="px-6 py-5 text-right">
        <button className="text-sc-on-surface-variant hover:text-sc-on-surface transition-colors">
          <Icon name="more_vert" />
        </button>
      </td>
    </tr>
  );
}
