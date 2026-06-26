import { FormEvent, useState } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { createPost } from '../lib/api';
import { FEED_REFRESH_EVENT } from '../lib/events';

export function AppShell() {
  const { signOut, sessionUser, username } = useAuth();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [posting, setPosting] = useState(false);
  const [error, setError] = useState('');

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPosting(true);
    setError('');
    try {
      await createPost({ title, content });
      setTitle('');
      setContent('');
      window.dispatchEvent(new Event(FEED_REFRESH_EVENT));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not create post');
    } finally {
      setPosting(false);
    }
  }

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

        <div className="compose-sidebar">
          <div className="section-tag">Create Post</div>
          <h3>Publish to feed</h3>
          <p className="hero-copy">This composer lives in the sidebar. The feed page stays read-only.</p>

          <form className="compose-form compose-form-sidebar" onSubmit={onSubmit}>
            <label>
              <span>Title</span>
              <input value={title} onChange={(event) => setTitle(event.target.value)} minLength={3} maxLength={30} required />
            </label>
            <label>
              <span>Content</span>
              <textarea value={content} onChange={(event) => setContent(event.target.value)} minLength={5} rows={5} required />
            </label>
            {error ? <div className="form-error">{error}</div> : null}
            <button type="submit" className="button" disabled={posting}>
              {posting ? 'Publishing...' : 'Create post'}
            </button>
          </form>
        </div>
      </aside>

      <main className="app-main">
        <Outlet />
      </main>
    </div>
  );
}
