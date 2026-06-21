import { createContext, useContext, useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { login, logout, refreshSession, registerApiHooks, registerUser } from '../lib/api';
import { clearAccessToken, readAccessToken, writeAccessToken } from '../lib/session';
import { decodeSessionUser } from '../lib/jwt';
import type { SessionUser } from '../types';

type AuthStatus = 'booting' | 'anonymous' | 'authenticated';

type RegisterInput = {
  name: string;
  username: string;
  email: string;
  password: string;
};

type LoginInput = {
  username: string;
  password: string;
};

type AuthContextValue = {
  status: AuthStatus;
  accessToken: string | null;
  sessionUser: SessionUser;
  isAuthenticated: boolean;
  signIn: (input: LoginInput) => Promise<void>;
  signUp: (input: RegisterInput) => Promise<void>;
  signOut: () => Promise<void>;
  bootstrap: () => Promise<void>;
  refreshAccessToken: () => Promise<boolean>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

async function fetchFreshAccessToken() {
  const response = await refreshSession();
  if (!response.access_token) {
    return null;
  }

  return response.access_token;
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>('booting');
  const [accessToken, setAccessToken] = useState<string | null>(readAccessToken());
  const [sessionUser, setSessionUser] = useState<SessionUser>(decodeSessionUser(readAccessToken()));

  useEffect(() => {
    registerApiHooks({
      getAccessToken: readAccessToken,
      setAccessToken: (token) => {
        writeAccessToken(token);
        setAccessToken(token);
        setSessionUser(decodeSessionUser(token));
        setStatus('authenticated');
      },
      clearAccessToken: () => {
        clearAccessToken();
        setAccessToken(null);
        setSessionUser({ userId: null, role: null });
        setStatus('anonymous');
      },
      refreshSession: async () => {
        try {
          const freshToken = await fetchFreshAccessToken();
          if (!freshToken) {
            clearAccessToken();
            setAccessToken(null);
            setSessionUser({ userId: null, role: null });
            setStatus('anonymous');
            return false;
          }

          writeAccessToken(freshToken);
          setAccessToken(freshToken);
          setSessionUser(decodeSessionUser(freshToken));
          setStatus('authenticated');
          return true;
        } catch {
          clearAccessToken();
          setAccessToken(null);
          setSessionUser({ userId: null, role: null });
          setStatus('anonymous');
          return false;
        }
      },
    });
  }, []);

  async function bootstrap() {
    const token = readAccessToken();
    if (token) {
      setAccessToken(token);
      setSessionUser(decodeSessionUser(token));
      setStatus('authenticated');
      return;
    }

    try {
      const freshToken = await fetchFreshAccessToken();
      if (!freshToken) {
        setStatus('anonymous');
        return;
      }

      writeAccessToken(freshToken);
      setAccessToken(freshToken);
      setSessionUser(decodeSessionUser(freshToken));
      setStatus('authenticated');
    } catch {
      clearAccessToken();
      setAccessToken(null);
      setSessionUser({ userId: null, role: null });
      setStatus('anonymous');
    }
  }

  useEffect(() => {
    void bootstrap();
  }, []);

  async function signIn(input: LoginInput) {
    const response = await login(input);
    if (!response.access_token) {
      throw new Error('Login succeeded but no access token was returned');
    }

    writeAccessToken(response.access_token);
    setAccessToken(response.access_token);
    setSessionUser(decodeSessionUser(response.access_token));
    setStatus('authenticated');
  }

  async function signUp(input: RegisterInput) {
    const response = await registerUser(input);
    if (!response.access_token) {
      throw new Error('Registration succeeded but no access token was returned');
    }

    writeAccessToken(response.access_token);
    setAccessToken(response.access_token);
    setSessionUser(decodeSessionUser(response.access_token));
    setStatus('authenticated');
  }

  async function signOut() {
    try {
      await logout();
    } finally {
      clearAccessToken();
      setAccessToken(null);
      setSessionUser({ userId: null, role: null });
      setStatus('anonymous');
    }
  }

  async function refreshAccessToken() {
    try {
      const freshToken = await fetchFreshAccessToken();
      if (!freshToken) {
        return false;
      }

      writeAccessToken(freshToken);
      setAccessToken(freshToken);
      setSessionUser(decodeSessionUser(freshToken));
      setStatus('authenticated');
      return true;
    } catch {
      clearAccessToken();
      setAccessToken(null);
      setSessionUser({ userId: null, role: null });
      setStatus('anonymous');
      return false;
    }
  }

  const value: AuthContextValue = {
    status,
    accessToken,
    sessionUser,
    isAuthenticated: status === 'authenticated',
    signIn,
    signUp,
    signOut,
    bootstrap,
    refreshAccessToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return value;
}
