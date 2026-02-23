/**
 * Axios instance with JWT interceptor and global error handling.
 * All API calls go through this configured client.
 */
import axios from "axios";
import router from "@/router";

const API_BASE = import.meta.env.VITE_API_BASE_URL || "";

const api = axios.create({
  baseURL: API_BASE,
  timeout: 15000,
  headers: { "Content-Type": "application/json" },
});

// ── Request interceptor: inject JWT ──
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// ── Response interceptor: handle 401 globally ──
api.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("access_token");
      router.push("/login");
    }
    return Promise.reject(error);
  }
);

export default api;
