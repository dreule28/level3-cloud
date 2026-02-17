<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "../api/auth";

const router = useRouter();

const username = ref("admin");
const password = ref("password");
const loading = ref(false);
const error = ref("");

async function onLogin() {
  error.value = "";
  loading.value = true;
  try {
    await login(username.value, password.value);
    router.push("/instances");
  } catch (e) {
    error.value = e?.body?.message || e.message || "Login failed";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div style="max-width: 360px; margin: 80px auto; font-family: system-ui;">
    <h2>Login</h2>

    <label>Username</label>
    <input v-model="username" style="width: 100%; padding: 8px; margin: 6px 0 12px;" />

    <label>Password</label>
    <input v-model="password" type="password" style="width: 100%; padding: 8px; margin: 6px 0 12px;" />

    <button @click="onLogin" :disabled="loading" style="width: 100%; padding: 10px;">
      {{ loading ? "Logging in..." : "Login" }}
    </button>

    <p v-if="error" style="color: #b00020; margin-top: 12px;">
      {{ error }}
    </p>
  </div>
</template>
