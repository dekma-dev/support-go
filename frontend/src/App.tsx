import { useEffect, useState } from "react";
import { Route, Routes, Navigate } from "react-router-dom";
import { AuthSession, readStoredSession } from "./api";
import { Layout } from "./components/Layout";
import { LoginPage } from "./pages/LoginPage";
import { DashboardHome } from "./pages/DashboardHome";
import { AgentDashboard } from "./pages/AgentDashboard";
import { TicketListPage } from "./pages/TicketListPage";
import { CreateTicketPage } from "./pages/CreateTicketPage";
import { TicketDetail } from "./TicketDetail";
import { UserManagement } from "./pages/UserManagement";
import { KnowledgeBase } from "./pages/KnowledgeBase";

export function App() {
  const [session, setSession] = useState<AuthSession | null>(() => readStoredSession());

  useEffect(() => {
    setSession(readStoredSession());
  }, []);

  if (!session) {
    return <LoginPage onLogin={setSession} />;
  }

  const isStaff = session.role === "agent" || session.role === "admin";

  return (
    <Layout session={session} onSessionChange={setSession}>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route
          path="/dashboard"
          element={isStaff ? <AgentDashboard session={session} /> : <DashboardHome />}
        />
        <Route path="/tickets" element={<TicketListPage session={session} />} />
        <Route path="/tickets/new" element={<CreateTicketPage />} />
        <Route path="/tickets/:id" element={<TicketDetail session={session} />} />
        <Route path="/users" element={<UserManagement session={session} />} />
        <Route path="/knowledge" element={<KnowledgeBase />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </Layout>
  );
}
