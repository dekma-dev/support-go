import { useState } from "react";
import { Icon } from "../components/Icon";
import { login, register, type AuthSession } from "../api";

export function LoginPage({ onLogin }: { onLogin: (s: AuthSession) => void }) {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const session =
        mode === "login"
          ? await login({ email, password })
          : await register({ email, password, role: "client" });
      onLogin(session);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Authentication failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex flex-col items-center justify-center p-6 relative"
      style={{
        backgroundColor: "#0b1326",
        backgroundImage:
          "radial-gradient(circle at 50% 50%, rgba(0,229,255,0.05) 0%, transparent 70%), linear-gradient(rgba(59,73,76,0.05) 1px, transparent 1px), linear-gradient(90deg, rgba(59,73,76,0.05) 1px, transparent 1px)",
        backgroundSize: "100% 100%, 40px 40px, 40px 40px",
      }}
    >
      {/* Atmospheric blurs */}
      <div className="fixed top-1/4 -right-20 w-96 h-96 bg-sc-primary-container/5 blur-[120px] rounded-full pointer-events-none" />
      <div className="fixed bottom-1/4 -left-20 w-96 h-96 bg-sc-secondary-container/5 blur-[120px] rounded-full pointer-events-none" />

      {/* Brand header */}
      <header className="mb-8 text-center relative z-10">
        <div className="flex items-center justify-center gap-3 mb-2">
          <Icon name="terminal" className="text-sc-primary-container text-3xl" />
          <h1 className="font-headline font-bold text-2xl tracking-widest text-sc-primary-container uppercase">
            SYNTHETIC COMMAND
          </h1>
        </div>
        <p className="font-label text-xs tracking-[0.25em] text-sc-on-surface-variant/70 uppercase">
          High-Tech Support Access Portal
        </p>
      </header>

      {/* Auth card */}
      <main className="w-full max-w-md relative z-10">
        {/* Corner accents */}
        <div className="absolute -top-1 -left-1 w-8 h-8 border-t-2 border-l-2 border-sc-primary-container/40 z-10" />
        <div className="absolute -bottom-1 -right-1 w-8 h-8 border-b-2 border-r-2 border-sc-primary-container/40 z-10" />

        <div className="bg-sc-surface-container shadow-2xl overflow-hidden relative border border-sc-outline-variant/10 rounded" style={{ boxShadow: "0 0 15px rgba(0,229,255,0.1), inset 0 0 1px rgba(195,245,255,0.3)", backdropFilter: "blur(20px)" }}>
          {/* Shimmer top edge */}
          <div className="h-px w-full" style={{ background: "linear-gradient(135deg, rgba(195,245,255,0.15) 0%, rgba(233,179,255,0.05) 100%)" }} />

          <div className="p-8 md:p-10">
            <div className="mb-10 text-center">
              <h2 className="font-headline font-bold text-3xl text-sc-on-surface mb-2 tracking-tight">
                SYSTEM AUTHENTICATION
              </h2>
              <div className="flex items-center justify-center gap-2">
                <span className="w-2 h-2 rounded-full bg-sc-primary-container animate-pulse" />
                <span className="font-label text-[10px] font-semibold text-sc-primary-container uppercase tracking-widest">
                  Secure Uplink Required
                </span>
              </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Email field */}
              <div className="space-y-2">
                <label className="font-label text-xs font-semibold text-sc-on-surface-variant uppercase tracking-wider block">
                  Employee ID
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Icon name="badge" className="text-sc-on-surface-variant text-sm" />
                  </div>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="agent@synth.cmd"
                    required
                    className="w-full bg-sc-surface-highest border-none rounded-sm py-3 pl-10 pr-4 text-sc-on-surface placeholder:text-sc-on-surface-variant/30 focus:ring-1 focus:ring-sc-primary focus:bg-sc-surface-bright transition-all outline-none"
                  />
                </div>
              </div>

              {/* Password field */}
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <label className="font-label text-xs font-semibold text-sc-on-surface-variant uppercase tracking-wider block">
                    Authorization Key
                  </label>
                  <span className="font-label text-[10px] font-semibold text-sc-primary-container/70 uppercase tracking-widest">
                    Forgot access code?
                  </span>
                </div>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Icon name="key" className="text-sc-on-surface-variant text-sm" />
                  </div>
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••••••"
                    required
                    className="w-full bg-sc-surface-highest border-none rounded-sm py-3 pl-10 pr-4 text-sc-on-surface placeholder:text-sc-on-surface-variant/30 focus:ring-1 focus:ring-sc-primary focus:bg-sc-surface-bright transition-all outline-none"
                  />
                </div>
              </div>

              {error && (
                <div className="text-sc-error text-xs font-label bg-sc-error/10 border border-sc-error/20 px-4 py-2 rounded-sm">
                  {error}
                </div>
              )}

              {/* Submit */}
              <div className="pt-4">
                <button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-sc-primary-fixed text-sc-on-primary-fixed font-headline font-bold py-4 rounded-sm tracking-widest flex items-center justify-center gap-2 shadow-[0_0_15px_rgba(0,229,255,0.2)] hover:shadow-[0_0_25px_rgba(0,229,255,0.4)] active:scale-[0.98] transition-all disabled:opacity-50"
                >
                  {loading
                    ? "AUTHENTICATING..."
                    : mode === "login"
                      ? "INITIATE COMMAND SEQUENCE"
                      : "REGISTER OPERATIVE"}
                  {!loading && <Icon name="login" className="text-lg" />}
                </button>
              </div>
            </form>

            {/* Toggle mode */}
            <div className="mt-8 pt-8 border-t border-sc-outline-variant/10 text-center">
              <p className="font-label text-xs text-sc-on-surface-variant">
                Unauthorized access is monitored.
              </p>
              <button
                onClick={() => setMode(mode === "login" ? "register" : "login")}
                className="inline-block mt-4 font-headline text-xs font-bold text-sc-on-surface hover:text-sc-primary-container transition-colors uppercase tracking-[0.2em]"
              >
                {mode === "login" ? "New Operative Registration" : "Existing Operative Login"}
              </button>
            </div>
          </div>

          {/* Bottom data bar */}
          <div className="bg-sc-surface-lowest px-4 py-2 flex justify-between items-center">
            <div className="flex gap-4">
              <div className="flex items-center gap-1">
                <span className="font-label text-[8px] text-sc-on-surface-variant/50 uppercase">LATENCY</span>
                <span className="font-headline text-[8px] text-sc-tertiary-fixed">12MS</span>
              </div>
              <div className="flex items-center gap-1">
                <span className="font-label text-[8px] text-sc-on-surface-variant/50 uppercase">NODE</span>
                <span className="font-headline text-[8px] text-sc-on-surface-variant">ALPHA-09</span>
              </div>
            </div>
            <div className="flex items-center gap-1">
              <Icon name="verified_user" className="text-[10px] text-sc-tertiary-fixed" />
              <span className="font-headline text-[8px] text-sc-on-surface-variant/80 uppercase">End-to-End Encryption Active</span>
            </div>
          </div>
        </div>
      </main>

      {/* Security badge */}
      <div className="mt-8 flex items-center gap-2 bg-sc-surface-low px-4 py-2 rounded-full border border-sc-outline-variant/10 shadow-sm relative z-10">
        <Icon name="shield" className="text-sc-primary-container text-sm" />
        <span className="font-label text-[10px] text-sc-on-surface-variant uppercase tracking-widest">
          Verified Tactical Terminal v4.2.0
        </span>
      </div>
    </div>
  );
}
