/**
 * Auth Store — Pinia
 * Manages JWT auth state, login/logout, and user info
 */
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { loginApi, logoutApi, getToken } from "@/api/auth";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(getToken());
  const loading = ref(false);
  const error = ref("");

  const isAuthenticated = computed(() => !!token.value);

  async function login(username, password) {
    loading.value = true;
    error.value = "";
    try {
      const data = await loginApi(username, password);
      token.value = data.access_token;
      return true;
    } catch (e) {
      error.value = e.response?.data?.message || e.message || "Login failed";
      return false;
    } finally {
      loading.value = false;
    }
  }

  function logout() {
    logoutApi();
    token.value = null;
  }

  function clearError() {
    error.value = "";
  }

  return { token, loading, error, isAuthenticated, login, logout, clearError };
});
