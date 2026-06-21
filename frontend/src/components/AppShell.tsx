import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export function AppShell() {
  const { signOut, sessionUser } = useAuth();

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand-block">
          <div className="brand-mark">B</div>
          <div>
            <div className="brand-name">Bozena Media</div>
            <div className="brand-subtitle">social microservices</div>
          </div>
        </div>

        <nav className="nav-stack">
          <NavLink to="/app" end className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Feed
          </NavLink>
          <NavLink to="/app/search" className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Search
          </NavLink>
          <NavLink
            to={sessionUser.userId ? `/app/users/${sessionUser.userId}` : '/app/users/1'}
            className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}
          >
            My profile
          </NavLink>
        </nav>

        <div className="sidebar-note">
          <p>Gateway-driven frontend</p>
          <span>Access token in session storage. Refresh token in HttpOnly cookie.</span>
        </div>

        <button type="button" className="button button-ghost" onClick={() => void signOut()}>
          Logout
        </button>
      </aside>

      <main className="app-main">
        <Outlet />
      </main>
    </div>
  );
}
