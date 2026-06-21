import { Link } from 'react-router-dom';

export function NotFoundPage() {
  return (
    <div className="screen-center">
      <div className="auth-card">
        <div className="eyebrow">404</div>
        <h1>Page not found</h1>
        <p className="hero-copy">The route does not exist in this frontend shell.</p>
        <Link className="button" to="/">
          Go home
        </Link>
      </div>
    </div>
  );
}
