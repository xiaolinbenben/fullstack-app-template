export type Envelope<T> = {
  success: boolean;
  message: string;
  data: T;
};

async function request<T>(path: string, init?: RequestInit): Promise<Envelope<T>> {
  const response = await fetch(path, { ...init, credentials: "same-origin" });
  const payload = (await response.json()) as Envelope<T>;
  if (!response.ok || !payload.success) {
    const error = new Error(payload.message || "请求失败") as Error & { status: number };
    error.status = response.status;
    throw error;
  }
  return payload;
}

export function fetchHealth(): Promise<Envelope<{ status: string }>> {
  return request("/healthz");
}

export function fetchWoodfish(): Promise<Envelope<{ count: number }>> {
  return request("/api/woodfish");
}

export function fetchSession(): Promise<Envelope<{ username: string }>> {
  return request("/api/admin/session");
}

export function login(username: string, password: string): Promise<Envelope<{ username: string }>> {
  return request("/api/admin/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
}

export function logout(): Promise<Envelope<Record<string, never>>> {
  return request("/api/admin/logout", { method: "POST" });
}
