import { http, setToken, clearToken } from "./http";

export async function login(username, password) {
  const data = await http.post("/auth/login", { username, password });
  // data = { access_token, token_type, expires_in }
  setToken(data.access_token);
  return data;
}

export function logout() {
  clearToken();
}
