<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { listInstances, createInstance, deleteInstance } from "../api/instances";
import { logout } from "../api/auth";

const router = useRouter();

const items = ref([]);
const loading = ref(false);
const error = ref("");

const newId = ref("");

async function load() {
  error.value = "";
  loading.value = true;
  try {
    items.value = await listInstances();
  } catch (e) {
    if (e.status === 401) {
      // token cleared in http.js already
      router.push("/login");
      return;
    }
    error.value = e?.body?.message || e.message || "Failed to load instances";
  } finally {
    loading.value = false;
  }
}

async function onCreate() {
  if (!newId.value.trim()) return;
  error.value = "";
  try {
    await createInstance(newId.value.trim());
    newId.value = "";
    await load();
  } catch (e) {
    error.value = e?.body?.message || e.message || "Create failed";
  }
}

async function onDelete(id) {
  if (!confirm(`Delete instance "${id}"? This cannot be undone.`)) return;
  error.value = "";
  try {
    await deleteInstance(id);
    await load();
  } catch (e) {
    error.value = e?.body?.message || e.message || "Delete failed";
  }
}

function onLogout() {
  logout();
  router.push("/login");
}

onMounted(load);
</script>

<template>
  <div class="instances-container">
    <div class="instances-header">
      <h2>Instances</h2>
      <button @click="onLogout">Logout</button>
    </div>

    <div class="create-row">
      <input
        v-model="newId"
        placeholder="new instance id (e.g. pg-demo)"
        class="create-input"
      />
      <button @click="onCreate">Create</button>
      <button @click="load" :disabled="loading">{{ loading ? "Loading..." : "Refresh" }}</button>
    </div>

    <p v-if="error" class="error-msg">{{ error }}</p>

    <table class="instances-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Status</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="it in items" :key="it.id">
          <td><router-link :to="`/instances/${it.id}`">{{ it.id }}</router-link></td>
          <td>{{ it.status }}</td>
          <td class="actions">
            <button @click="onDelete(it.id)">Delete</button>
          </td>
        </tr>
        <tr v-if="!loading && items.length === 0">
          <td colspan="3" class="empty-msg">No instances.</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.instances-container {
  max-width: 820px;
  margin: 40px auto;
  font-family: system-ui;
}
.instances-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.create-row {
  display: flex;
  gap: 8px;
  margin: 12px 0 18px;
}
.create-input {
  flex: 1;
  padding: 8px;
}
.error-msg {
  color: #b00020;
}
.instances-table {
  width: 100%;
  border-collapse: collapse;
}
.instances-table th {
  text-align: left;
  border-bottom: 1px solid #ddd;
  padding: 8px;
}
.instances-table td {
  padding: 8px;
  border-bottom: 1px solid #eee;
}
.instances-table .actions {
  text-align: right;
}
.empty-msg {
  padding: 12px;
  color: #666;
}
</style>
