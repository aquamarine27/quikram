import { createContext, useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { login as apiLogin, refresh as apiRefresh, logout as apiLogout, HAS_SESSION_KEY } from "../api/auth";
import type { LoginData } from "../api/auth";
import { setAccessToken as setClientToken } from "../api/client";

// ─── Types ───
interface User {
  id: string;
  email: string;
  name: string;
  plan: string;
  uploads_this_month: number;
  uploads_reset_at: string | null;
  created_at: string;
  updated_at: string;
}

interface AuthContextType {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (data: LoginData) => Promise<void>;
  logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | null>(null);

// ─── Provider ───
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const disposed = useRef(false);

  // Restore session on mount (skip if no stored flag)
  useEffect(() => {
    const hasSession = localStorage.getItem(HAS_SESSION_KEY) === "true";
    if (!hasSession) {
      setIsLoading(false);
      return;
    }
    (async () => {
      try {
        const { data } = await apiRefresh();
        setUser(data.user);
        setAccessToken(data.access_token);
        setClientToken(data.access_token);
      } catch {
        setUser(null);
        setAccessToken(null);
        localStorage.removeItem(HAS_SESSION_KEY);
      } finally {
        if (!disposed.current) setIsLoading(false);
      }
    })();
    return () => { disposed.current = true; };
  }, []);

  const login = useCallback(async (data: LoginData) => {
    const { data: res } = await apiLogin(data);
    setUser(res.user);
    setAccessToken(res.access_token);
    setClientToken(res.access_token);
    localStorage.setItem(HAS_SESSION_KEY, "true");
  }, []);

  const logout = useCallback(async () => {
    try {
      await apiLogout();
    } finally {
      setUser(null);
      setAccessToken(null);
      setClientToken(null);
      localStorage.removeItem(HAS_SESSION_KEY);
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, accessToken, isAuthenticated: !!user, isLoading, login, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}
