import { Link } from 'react-router-dom';

export function LandingPage() {
  return (
    <div className="landing-page">
      <div className="landing-orb landing-orb-one" />
      <div className="landing-orb landing-orb-two" />

      <header className="landing-header">
        <div className="brand-block">
          <div className="brand-mark">B</div>
          <div>
            <div className="brand-name">Bozena Media</div>
            <div className="brand-subtitle">microservice social frontend</div>
          </div>
        </div>
        <div className="landing-actions">
          <Link className="button button-ghost" to="/login">
            Login
          </Link>
          <Link className="button" to="/register">
            Create account
          </Link>
        </div>
      </header>

      <main className="landing-hero">
        <div className="hero-copy-block">
          <div className="eyebrow">React frontend over the existing Go gateway</div>
          <h1>One app for posts, follows, likes, and auth.</h1>
          <p>
            This interface talks to the API Gateway exactly as it exists today, using the browser client flow with
            cookies for refresh tokens and bearer tokens for protected routes.
          </p>
          <div className="landing-actions">
            <Link className="button" to="/register">
              Get started
            </Link>
            <Link className="button button-soft" to="/app">
              Open dashboard
            </Link>
          </div>
        </div>

        <div className="landing-panel">
          <div className="glass-card">
            <span className="section-tag">Current coverage</span>
            <ul className="feature-list">
              <li>Auth: register, login, refresh, logout, search</li>
              <li>Posts: list, create, edit, delete, by-user</li>
              <li>Social: likes, followers, followings</li>
              <li>Cursor pagination where the backend already supports it</li>
            </ul>
          </div>
        </div>
      </main>
    </div>
  );
}
