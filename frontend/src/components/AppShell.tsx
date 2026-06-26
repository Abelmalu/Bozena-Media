import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export function AppShell() {
  const { signOut, sessionUser, username } = useAuth();

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand-block">
          <div className="brand-mark">B</div>
          <div>
            <div className="brand-name">Bozena Media</div>
            <div className="brand-subtitle">microservice social app</div>
          </div>
        </div>

        <div className="profile-card">
          <div className="profile-label">Signed in as</div>
          <div className="profile-value">{username ?? 'Unknown user'}</div>
          <div className="profile-subvalue">{sessionUser.userId ? `User #${sessionUser.userId}` : 'ID unavailable'}</div>
          <div className="profile-subvalue">{sessionUser.role ?? 'role unavailable'}</div>
        </div>

        <nav className="nav-stack">
          <NavLink to="/app/feed" className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Feed
          </NavLink>
          <NavLink to="/app/create-post" className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Create Post
          </NavLink>
          <NavLink to="/app/profile" className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Profile
          </NavLink>
          <NavLink to="/app/search" className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Search
          </NavLink>
        </nav>

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
