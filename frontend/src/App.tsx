import { useEffect, useState } from "react";
import { Route, Routes, Navigate } from "react-router-dom";
import { AuthSession, readStoredSession } from "./api";
import { Layout } from "./components/Layout";
import { LoginPage } from "./pages/LoginPage";
import { DashboardHome } from "./pages/DashboardHome";
import { TicketListPage } from "./pages/TicketListPage";
import { TicketDetail } from "./TicketDetail";

export function App() {
  const [session, setSession] = useState<AuthSession | null>(() => readStoredSession());

  useEffect(() => {
    setSession(readStoredSession());
  }, []);

  if (!session) {
    return <LoginPage onLogin={setSession} />;
  }

  return (
    <Layout session={session} onSessionChange={setSession}>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/dashboard" element={<DashboardHome />} />
        <Route path="/tickets" element={<TicketListPage session={session} />} />
        <Route path="/tickets/:id" element={<TicketDetail session={session} />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Layout>
  );
}
