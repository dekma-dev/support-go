import { Link, Route, Routes } from "react-router-dom";

function HomePage() {
  return (
    <section className="card">
      <h2>Support-Go UI</h2>
      <p>Фундамент frontend на React + TypeScript + Vite готов.</p>
    </section>
  );
}

function TicketsPage() {
  return (
    <section className="card">
      <h2>Tickets</h2>
      <p>Здесь будет список тикетов, фильтры и действия агента.</p>
    </section>
  );
}

export function App() {
  return (
    <main className="layout">
      <header className="header">
        <h1>Support-Go</h1>
        <nav className="nav">
          <Link to="/">Home</Link>
          <Link to="/tickets">Tickets</Link>
        </nav>
      </header>

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/tickets" element={<TicketsPage />} />
      </Routes>
    </main>
  );
}

