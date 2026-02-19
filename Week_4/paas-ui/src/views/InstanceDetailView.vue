<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getInstance } from "../api/instances";

const route = useRoute();
const router = useRouter();

const instance = ref(null);
const loading = ref(false);
const error = ref("");
const showPassword = ref(false);

async function load() {
  loading.value = true;
  error.value = "";
  try {
    instance.value = await getInstance(route.params.id);
  } catch (e) {
    if (e.status === 401) {
      router.push("/login");
      return;
    }
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}

function copy(text) {
  navigator.clipboard.writeText(text);
}

onMounted(load);
</script>

<template>
  <div class="detail-container">
    <button @click="$router.push('/instances')">← Back</button>

    <h2>Instance: {{ route.params.id }}</h2>

    <p v-if="loading">Loading...</p>
    <p v-if="error" class="error-msg">{{ error }}</p>

    <div v-if="instance">
      <p><strong>Status:</strong> {{ instance.status }}</p>

      <div v-if="instance.connection" class="connection-section">
        <h3>Connection</h3>
        <p><strong>Host:</strong> {{ instance.connection.host }}</p>
        <p><strong>Port:</strong> {{ instance.connection.port }}</p>
        <p><strong>Database:</strong> {{ instance.connection.database }}</p>
        <p><strong>User:</strong> {{ instance.connection.user }}</p>
        <p>
          <strong>Password:</strong>
          <code v-if="showPassword">{{ instance.connection.password }}</code>
          <code v-else>********</code>
          <button @click="showPassword = !showPassword">{{ showPassword ? 'Hide' : 'Show' }}</button>
          <button @click="copy(instance.connection.password)">Copy</button>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-container {
  max-width: 800px;
  margin: 40px auto;
  font-family: system-ui;
}
.error-msg {
  color: #b00020;
}
.connection-section {
  margin-top: 20px;
}
.connection-section code {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  margin: 0 4px;
}
.connection-section button {
  margin-left: 4px;
}
</style>
