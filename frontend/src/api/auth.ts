// ─── App state ───
export const HAS_SESSION_KEY = "quikram_session";

// ─── Auth type ───
export interface LoginData {
  email: string;
  password: string;
}

// ─── Register payload ───
export interface RegisterData {
  email: string;
  password: string;
  name: string;
}

// ─── Login response ───
export interface LoginResponse {
  access_token: string;
  user: {
    id: string;
    email: string;
    name: string;
    plan: string;
  };
}

const API_BASE = import.meta.env.VITE_API_URL || "/api/v1";
const JSON_HEADERS = { "Content-Type": "application/json" };

// ─── POST /auth/... ───
async function authPost<T>(path: string, body: object): Promise<{ data: T }> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: JSON_HEADERS,
    credentials: "include",
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }
  return { data: await res.json() };
}

export function login(data: LoginData) {
  return authPost<LoginResponse>("/auth/login", data);
}

export function register(data: RegisterData) {
  return authPost<LoginResponse>("/auth/register", data);
}

export function refresh() {
  return authPost<LoginResponse>("/auth/refresh", {});
}

export function logout() {
  return authPost<{ message: string }>("/auth/logout", {});
}

export function getProfile() {
  return authPost<LoginResponse["user"]>("/auth/profile", {});
}
