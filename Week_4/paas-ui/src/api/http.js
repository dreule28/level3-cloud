const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

export function getToken() {
  return localStorage.getItem("access_token");
}

export function setToken(token) {
  localStorage.setItem("access_token", token);
}

export function clearToken() {
  localStorage.removeItem("access_token");
}

async function request(path, options = {}) {
  const token = getToken();

  const headers = {
    ...(options.headers || {}),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  // Auto-handle 401 globally
  if (res.status === 401) {
    clearToken();
    // Let caller handle navigation; throw a structured error
    const body = await safeJson(res);
    const err = new Error(body?.message || "Unauthorized");
    err.status = 401;
    err.body = body;
    throw err;
  }

  if (!res.ok) {
    const body = await safeJson(res);
    const err = new Error(body?.message || `HTTP ${res.status}`);
    err.status = res.status;
    err.body = body;
    throw err;
  }

  // Some endpoints may return empty body
  return safeJson(res);
}

async function safeJson(res) {
  const text = await res.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return { raw: text };
  }
}

export const http = {
  get: (path) => request(path),
  post: (path, body) =>
    request(path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    }),
  del: (path) =>
    request(path, {
      method: "DELETE",
    }),
};
