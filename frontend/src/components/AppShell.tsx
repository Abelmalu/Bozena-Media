import { useEffect, useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { connectNotificationStream } from '../lib/sse';
import { NOTIFICATION_RECEIVED_EVENT } from '../lib/events';

export function AppShell() {
  const { signOut, sessionUser, username, isAuthenticated } = useAuth();
  const [toast, setToast] = useState('');

  useEffect(() => {
    if (!isAuthenticated) return;

    // Connect to SSE with Bearer token for real-time notifications
    const controller = connectNotificationStream((data) => {
      setToast(data);
      // Dispatch custom event to notify other components (e.g. NotificationsPage)
      window.dispatchEvent(new CustomEvent(NOTIFICATION_RECEIVED_EVENT, { detail: data }));
      setTimeout(() => setToast(''), 5000);
    });

    return () => controller.abort();
  }, [isAuthenticated]);

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
          <NavLink to="/app/notifications" className={({ isActive }) => (isActive ? 'nav-link active' : 'nav-link')}>
            Notifications
          </NavLink>
        </nav>

        <button type="button" className="button button-ghost" onClick={() => void signOut()}>
          Logout
        </button>
      </aside>

      <main className="app-main">
        {toast && (
          <div className="toast" style={{ position: 'fixed', top: 16, right: 16, background: '#0ea5e9', color: '#fff', padding: '12px 24px', borderRadius: '8px', zIndex: 9999, boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}>
            <strong>New Notification:</strong> {toast}
          </div>
        )}
        <Outlet />
      </main>
    </div>
  );
}
