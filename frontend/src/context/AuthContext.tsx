import { createContext, useContext, useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { clearAccessToken, clearUsername, readAccessToken, readUsername, writeAccessToken, writeUsername } from '../lib/session';
import { decodeSessionUser } from '../lib/jwt';
import { login, logout, refreshSession, registerAuthHooks, registerUser } from '../lib/api';
import type { SessionUser } from '../types';

type AuthState = 'booting' | 'anonymous' | 'authenticated';

type AuthContextValue = {
  state: AuthState;
  accessToken: string | null;
  sessionUser: SessionUser;
  username: string | null;
  isAuthenticated: boolean;
  signIn: (input: { username: string; password: string }) => Promise<void>;
  signUp: (input: { name: string; username: string; email: string; password: string }) => Promise<void>;
  signOut: () => Promise<void>;
  refreshAccessToken: () => Promise<boolean>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>('booting');
  const [accessToken, setAccessTokenState] = useState<string | null>(readAccessToken());
  const [sessionUser, setSessionUser] = useState<SessionUser>(decodeSessionUser(readAccessToken()));
  const [username, setUsername] = useState<string | null>(readUsername());

  useEffect(() => {
    registerAuthHooks({
      getAccessToken: readAccessToken,
      setAccessToken: (token) => {
        writeAccessToken(token);
        setAccessTokenState(token);
        setSessionUser(decodeSessionUser(token));
        setState('authenticated');
      },
      clearAccessToken: () => {
        clearAccessToken();
        clearUsername();
        setAccessTokenState(null);
        setUsername(null);
        setSessionUser({ userId: null, role: null });
        setState('anonymous');
      },
      refreshSession: async () => {
        try {
          const response = await refreshSession();
          if (!response.access_token) {
            clearAccessToken();
            clearUsername();
            setAccessTokenState(null);
            setUsername(null);
            setSessionUser({ userId: null, role: null });
            setState('anonymous');
            return false;
          }

          writeAccessToken(response.access_token);
          if (response.username) {
            writeUsername(response.username);
            setUsername(response.username);
          }
          setAccessTokenState(response.access_token);
          setSessionUser(decodeSessionUser(response.access_token));
          setState('authenticated');
          return true;
        } catch {
          clearAccessToken();
          clearUsername();
          setAccessTokenState(null);
          setUsername(null);
          setSessionUser({ userId: null, role: null });
          setState('anonymous');
          return false;
        }
      },
    });
  }, []);

  useEffect(() => {
    async function bootstrap() {
      const token = readAccessToken();
      if (token) {
        setAccessTokenState(token);
        setSessionUser(decodeSessionUser(token));
        setState('authenticated');
        return;
      }

      try {
        const response = await refreshSession();
        if (response.access_token) {
          writeAccessToken(response.access_token);
          if (response.username) {
            writeUsername(response.username);
            setUsername(response.username);
          }
          setAccessTokenState(response.access_token);
          setSessionUser(decodeSessionUser(response.access_token));
          setState('authenticated');
        } else {
          setState('anonymous');
        }
      } catch {
        clearAccessToken();
        clearUsername();
        setAccessTokenState(null);
        setUsername(null);
        setSessionUser({ userId: null, role: null });
        setState('anonymous');
      }
    }

    void bootstrap();
  }, []);

  async function signIn(input: { username: string; password: string }) {
    const response = await login(input);
    if (!response.access_token) {
      throw new Error('Login succeeded without an access token');
    }

    writeAccessToken(response.access_token);
    if (response.username) {
      writeUsername(response.username);
      setUsername(response.username);
    }
    setAccessTokenState(response.access_token);
    setSessionUser(decodeSessionUser(response.access_token));
    setState('authenticated');
  }

  async function signUp(input: { name: string; username: string; email: string; password: string }) {
    const response = await registerUser(input);
    if (!response.access_token) {
      throw new Error('Registration succeeded without an access token');
    }

    writeAccessToken(response.access_token);
    if (response.username) {
      writeUsername(response.username);
      setUsername(response.username);
    }
    setAccessTokenState(response.access_token);
    setSessionUser(decodeSessionUser(response.access_token));
    setState('authenticated');
  }

  async function signOut() {
    try {
      await logout();
    } finally {
      clearAccessToken();
      clearUsername();
      setAccessTokenState(null);
      setUsername(null);
      setSessionUser({ userId: null, role: null });
      setState('anonymous');
    }
  }

  async function refreshAccessToken() {
    try {
      const response = await refreshSession();
      if (!response.access_token) {
        return false;
      }

      writeAccessToken(response.access_token);
      if (response.username) {
        writeUsername(response.username);
        setUsername(response.username);
      }
      setAccessTokenState(response.access_token);
      setSessionUser(decodeSessionUser(response.access_token));
      setState('authenticated');
      return true;
    } catch {
      clearAccessToken();
      clearUsername();
      setAccessTokenState(null);
      setUsername(null);
      setSessionUser({ userId: null, role: null });
      setState('anonymous');
      return false;
    }
  }

  const value: AuthContextValue = {
    state,
    accessToken,
    sessionUser,
    username,
    isAuthenticated: state === 'authenticated',
    signIn,
    signUp,
    signOut,
    refreshAccessToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error('useAuth must be inside AuthProvider');
  }
  return value;
}
