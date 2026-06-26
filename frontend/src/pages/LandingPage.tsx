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
            <div className="brand-subtitle">React frontend for the Go backend</div>
          </div>
        </div>

        <div className="landing-actions">
          <Link to="/login" className="button button-ghost">
            Login
          </Link>
          <Link to="/register" className="button">
            Create account
          </Link>
        </div>
      </header>

      <main className="landing-hero">
        <div>
          <div className="eyebrow">Gateway driven</div>
          <h1>Feed-first social app with profile, search, and actions.</h1>
          <p className="hero-copy">
            The frontend talks only to the API Gateway and uses the current backend contracts for auth, feed, search,
            followers, followings, likes, and post creation.
          </p>

          <div className="landing-actions">
            <Link to="/register" className="button">
              Start here
            </Link>
            <Link to="/login" className="button button-soft">
              Sign in
            </Link>
          </div>
        </div>

        <section className="glass-card">
          <div className="section-tag">Covered routes</div>
          <ul className="feature-list">
            <li>Auth: register, login, refresh, logout, search</li>
            <li>Feed: `/api/feed/` as the post-login default</li>
            <li>Profile: follower/following counts and pagination</li>
            <li>Search: cursor-based user lookup</li>
          </ul>
        </section>
      </main>
    </div>
  );
}
