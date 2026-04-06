import { TopNav } from "./TopNav";
import { Sidebar } from "./Sidebar";
import type { AuthSession } from "../api";

export function Layout({
  session,
  onSessionChange,
  children,
}: {
  session: AuthSession | null;
  onSessionChange: (s: AuthSession | null) => void;
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen w-full overflow-hidden">
      <Sidebar session={session} onSessionChange={onSessionChange} />
      <div className="flex-1 flex flex-col min-w-0">
        <TopNav session={session} />
        <main className="flex-1 overflow-y-auto bg-sc-surface-lowest">
          {children}
        </main>
      </div>
    </div>
  );
}
