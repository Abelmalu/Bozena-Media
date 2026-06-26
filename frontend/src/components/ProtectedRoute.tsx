import { Navigate, useLocation } from 'react-router-dom';
import type { ReactNode } from 'react';
import { useAuth } from '../context/AuthContext';

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { state, isAuthenticated } = useAuth();
  const location = useLocation();

  if (state === 'booting') {
    return (
      <div className="screen-center">
        <div className="loader-card">
          <div className="loader" />
          <p>Loading session...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />;
  }

  return <>{children}</>;
}
