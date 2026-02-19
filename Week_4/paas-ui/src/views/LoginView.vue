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
  <div class="login-container">
    <h2>Login</h2>

    <label>Username</label>
    <input v-model="username" class="login-input" />

    <label>Password</label>
    <input v-model="password" type="password" class="login-input" />

    <button @click="onLogin" :disabled="loading" class="login-button">
      {{ loading ? "Logging in..." : "Login" }}
    </button>

    <p v-if="error" class="login-error">
      {{ error }}
    </p>
  </div>
</template>

<style scoped>
.login-container {
  max-width: 360px;
  margin: 80px auto;
  font-family: system-ui;
}
.login-input {
  width: 100%;
  padding: 8px;
  margin: 6px 0 12px;
  box-sizing: border-box;
}
.login-button {
  width: 100%;
  padding: 10px;
  cursor: pointer;
}
.login-error {
  color: #b00020;
  margin-top: 12px;
}
</style>
