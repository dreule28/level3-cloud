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
  <div style="max-width: 820px; margin: 40px auto; font-family: system-ui;">
    <div style="display:flex; justify-content: space-between; align-items: center;">
      <h2>Instances</h2>
      <button @click="onLogout">Logout</button>
    </div>

    <div style="display:flex; gap: 8px; margin: 12px 0 18px;">
      <input
        v-model="newId"
        placeholder="new instance id (e.g. pg-demo)"
        style="flex:1; padding: 8px;"
      />
      <button @click="onCreate">Create</button>
      <button @click="load" :disabled="loading">{{ loading ? "Loading..." : "Refresh" }}</button>
    </div>

    <p v-if="error" style="color:#b00020;">{{ error }}</p>

    <table style="width:100%; border-collapse: collapse;">
      <thead>
        <tr>
          <th style="text-align:left; border-bottom: 1px solid #ddd; padding: 8px;">ID</th>
          <th style="text-align:left; border-bottom: 1px solid #ddd; padding: 8px;">Status</th>
          <th style="border-bottom: 1px solid #ddd; padding: 8px;"></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="it in items" :key="it.id">
          <td style="padding: 8px; border-bottom: 1px solid #eee;"><router-link :to="`/instances/${it.id}`">{{ it.id }}</router-link></td>
          <td style="padding: 8px; border-bottom: 1px solid #eee;">{{ it.status }}</td>
          <td style="padding: 8px; border-bottom: 1px solid #eee; text-align:right;">
            <button @click="onDelete(it.id)">Delete</button>
          </td>
        </tr>
        <tr v-if="!loading && items.length === 0">
          <td colspan="3" style="padding: 12px; color:#666;">No instances.</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
