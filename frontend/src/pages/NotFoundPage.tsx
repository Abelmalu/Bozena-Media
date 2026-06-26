import { Link } from 'react-router-dom';

export function NotFoundPage() {
  return (
    <div className="screen-center">
      <div className="auth-card">
        <div className="eyebrow">404</div>
        <h1>Not found</h1>
        <p className="hero-copy">That route does not exist.</p>
        <Link className="button" to="/">
          Go home
        </Link>
      </div>
    </div>
  );
}
