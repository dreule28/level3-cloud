/**
 * Auth API — login / logout / token management
 */
import api from "./axios";

export async function loginApi(username, password) {
  const { data } = await api.post("/auth/login", { username, password });
  localStorage.setItem("access_token", data.access_token);
  return data;
}

export function logoutApi() {
  localStorage.removeItem("access_token");
}

export function getToken() {
  return localStorage.getItem("access_token");
}
