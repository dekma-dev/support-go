import { Icon } from "../components/Icon";

const articles = [
  {
    id: "PROTO-NX-99",
    title: "Core Firewall Bypass Mitigation",
    description: "Standard operating procedure for isolating anomalous traffic patterns within the Tier 1 backbone. Includes specific routing table injection sequences.",
    tags: [{ label: "Network", color: "text-sc-secondary" }, { label: "Priority-High", color: "text-sc-error" }, { label: "L3-Access", color: "text-sc-on-surface-variant" }],
    author: "S. Vance",
    idColor: "bg-sc-primary/10 text-sc-primary-fixed-dim border-sc-primary/20",
  },
  {
    id: "DOC-SF-442",
    title: "Quantum DB Sync Latency",
    description: "Troubleshooting steps for asynchronous replication lag across distributed ledger nodes. Focuses on shard re-balancing and priority queueing.",
    tags: [{ label: "Software", color: "text-sc-tertiary" }, { label: "Critical-Fix", color: "text-sc-primary" }, { label: "Admin", color: "text-sc-on-surface-variant" }],
    author: "K. Chen",
    idColor: "bg-sc-tertiary/10 text-sc-tertiary border-sc-tertiary/20",
  },
  {
    id: "TEM-CX-01",
    title: "Priority Escalation Template",
    description: "Pre-formatted communication block for VIP client service interruptions. Contains automated variable fields for incident depth and estimated ETR.",
    tags: [{ label: "Templates", color: "text-sc-secondary" }, { label: "External", color: "text-sc-on-surface-variant" }],
    author: "M. Ross",
    idColor: "bg-sc-secondary/10 text-sc-secondary border-sc-secondary/20",
  },
  {
    id: "PROTO-HW-81",
    title: "Server Blade Hot-Swap Guide",
    description: "Physical maintenance protocol for Gen-9 server arrays. Critical safety steps to prevent static discharge and filesystem corruption during swap.",
    tags: [{ label: "Hardware", color: "text-sc-primary-fixed" }, { label: "Maintenance", color: "text-sc-on-surface-variant" }],
    author: "J. Diaz",
    idColor: "bg-sc-primary/10 text-sc-primary border-sc-primary/20",
  },
  {
    id: "SEC-AUT-09",
    title: "MFA Bypass Investigation",
    description: "Forensic analysis steps for identifying session hijacking through stolen authentication tokens. Includes log analysis regex patterns.",
    tags: [{ label: "Security", color: "text-sc-error" }, { label: "Incident", color: "text-sc-primary" }],
    author: "B. Kale",
    idColor: "bg-sc-error/10 text-sc-error border-sc-error/20",
  },
];

const categories = ["Network", "Software", "Hardware", "Security"];

export function KnowledgeBase() {
  return (
    <div className="min-h-full" style={{ backgroundImage: "radial-gradient(circle at 2px 2px, rgba(0,229,255,0.05) 1px, transparent 0)", backgroundSize: "32px 32px" }}>
      {/* Header */}
      <header className="p-8 lg:p-12 pb-0">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12">
          <div className="space-y-1">
            <div className="flex items-center gap-2 text-sc-primary text-[10px] font-bold tracking-[0.3em] uppercase">
              <span className="w-8 h-px bg-sc-primary" />
              Centralized Repository
            </div>
            <h1 className="font-headline text-5xl md:text-7xl font-bold tracking-tighter text-sc-on-surface">
              ARCHIVE <span className="text-sc-primary/40">&amp;</span> PROTOCOLS
            </h1>
            <p className="text-sc-on-surface-variant max-w-xl text-sm leading-relaxed mt-4">
              System-wide intelligence database. Access verified resolution templates, technical schematics, and operational workflows.
            </p>
          </div>
          <button className="bg-sc-primary-fixed text-sc-on-primary-fixed px-6 py-3 rounded-sm font-headline font-bold text-xs tracking-widest uppercase flex items-center gap-2 shadow-[0_0_15px_rgba(0,218,243,0.3)] hover:brightness-110 active:scale-95 transition-all shrink-0">
            <Icon name="add" className="text-sm" />
            CREATE NEW DOCUMENTATION
          </button>
        </div>

        {/* Search + filters */}
        <div className="glass-card p-2 rounded-md flex flex-wrap items-center gap-4 border border-sc-outline-variant/10">
          <div className="flex-1 min-w-[280px] relative">
            <Icon name="search" className="absolute left-4 top-1/2 -translate-y-1/2 text-sc-primary" />
            <input
              type="text"
              className="w-full bg-sc-surface-lowest border-none py-4 pl-12 pr-4 text-xs font-bold tracking-widest text-sc-on-surface placeholder:text-sc-on-surface-variant/40 focus:ring-1 focus:ring-sc-primary-container uppercase"
              placeholder="FILTER BY KEYWORD, TAG, OR PROTOCOL ID..."
            />
          </div>
          <div className="flex gap-2 p-2">
            {categories.map((cat, i) => (
              <button
                key={cat}
                className={`px-4 py-2 rounded-sm text-[10px] font-bold tracking-widest uppercase border transition-colors ${
                  i === 0
                    ? "bg-sc-primary/10 text-sc-primary border-sc-primary/20"
                    : "bg-sc-surface-container text-sc-on-surface-variant border-sc-outline-variant/20 hover:bg-sc-surface-highest"
                }`}
              >
                {cat}
              </button>
            ))}
          </div>
        </div>
      </header>

      {/* Cards grid */}
      <section className="p-8 lg:p-12 pt-12 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
        {articles.map((a) => (
          <article key={a.id} className="group relative flex flex-col glass-card p-6 border border-sc-outline-variant/10 hover:border-sc-primary/30 transition-all duration-300">
            <div className="flex justify-between items-start mb-6">
              <div className={`px-3 py-1 text-[9px] font-black tracking-widest border uppercase ${a.idColor}`}>
                {a.id}
              </div>
              <Icon name="open_in_new" className="text-sc-on-surface-variant/40 group-hover:text-sc-primary transition-colors cursor-pointer" />
            </div>
            <h3 className="font-headline text-xl font-bold text-sc-on-surface mb-3 group-hover:text-sc-primary transition-colors">{a.title}</h3>
            <p className="text-sc-on-surface-variant text-xs leading-relaxed mb-6 line-clamp-3">{a.description}</p>
            <div className="flex flex-wrap gap-2 mt-auto">
              {a.tags.map((tag) => (
                <span key={tag.label} className={`px-2 py-1 bg-sc-surface-highest text-[9px] font-bold uppercase tracking-tighter ${tag.color}`}>
                  {tag.label}
                </span>
              ))}
            </div>
            <div className="mt-6 pt-6 border-t border-sc-outline-variant/10 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="w-6 h-6 rounded-full bg-sc-surface-highest flex items-center justify-center text-[10px] font-bold text-sc-primary">
                  {a.author[0]}
                </div>
                <span className="text-[10px] font-bold text-sc-on-surface-variant/60 uppercase">Authored by: {a.author}</span>
              </div>
              <button className="text-[10px] font-bold text-sc-primary uppercase tracking-widest hover:underline">Quick Preview</button>
            </div>
          </article>
        ))}

        {/* Add new card */}
        <article className="group flex flex-col items-center justify-center border-2 border-dashed border-sc-outline-variant/20 rounded-md p-12 hover:border-sc-primary/40 hover:bg-sc-primary/5 transition-all cursor-pointer">
          <div className="w-16 h-16 rounded-full bg-sc-surface-container flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <Icon name="add_circle" className="text-3xl text-sc-on-surface-variant/40 group-hover:text-sc-primary" />
          </div>
          <span className="font-headline text-sm font-bold tracking-widest text-sc-on-surface-variant/60 uppercase group-hover:text-sc-primary">Add New Protocol</span>
          <span className="text-[10px] text-sc-on-surface-variant/40 uppercase mt-2">Level 3 Clearance Required</span>
        </article>
      </section>

      {/* Floating stats bar */}
      <div className="fixed bottom-6 right-6 hidden xl:flex gap-6 items-center glass-card px-6 py-4 rounded-full border border-sc-primary/20 shadow-[0_0_30px_rgba(0,218,243,0.15)] z-40">
        <div className="flex items-center gap-3 pr-6 border-r border-sc-outline-variant/20">
          <div className="w-2 h-2 rounded-full bg-sc-tertiary animate-pulse" />
          <div className="text-[10px] font-bold tracking-widest uppercase">
            <span className="text-sc-on-surface-variant/50">Database:</span> <span className="text-sc-on-surface">ONLINE</span>
          </div>
        </div>
        <div className="flex items-center gap-3 pr-6 border-r border-sc-outline-variant/20">
          <Icon name="description" className="text-sm text-sc-primary" />
          <div className="text-[10px] font-bold tracking-widest uppercase">
            <span className="text-sc-on-surface-variant/50">Docs:</span> <span className="text-sc-on-surface">{articles.length}</span>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <Icon name="update" className="text-sm text-sc-secondary" />
          <div className="text-[10px] font-bold tracking-widest uppercase text-sc-on-surface">
            Last Sync: <span className="text-sc-primary">2M AGO</span>
          </div>
        </div>
      </div>
    </div>
  );
}
