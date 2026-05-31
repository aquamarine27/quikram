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
    uploads_this_month: number;
    uploads_reset_at: string | null;
    created_at: string;
    updated_at: string;
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

// ─── POST /users/me/change-password ───
export interface ChangePasswordData {
  old_password: string;
  new_password: string;
}

export async function changePassword(data: ChangePasswordData): Promise<{ data: { message: string } }> {
  const { default: client } = await import("./client");
  const res = await client.post<{ message: string }>("/users/me/change-password", data);
  return { data: res.data };
}

// ─── POST /users/me/change-plan ───
export async function changePlan(plan: string): Promise<LoginResponse["user"]> {
  const { default: client } = await import("./client");
  const res = await client.post<LoginResponse["user"]>("/users/me/change-plan", { plan });
  return res.data;
}

// ─── GET /analytics/me ───
export interface AnalyticsData {
  subjects_count: number;
  total_tests: number;
  average_score: number;
}

export async function getAnalytics(): Promise<{ data: AnalyticsData }> {
  const { default: client } = await import("./client");
  const res = await client.get<AnalyticsData>("/analytics/me");
  return { data: res.data };
}
