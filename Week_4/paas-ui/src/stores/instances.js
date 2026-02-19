/**
 * Instances Store — Pinia
 * Manages instance list, details, CRUD state
 */
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  listInstances,
  getInstance,
  createInstance,
  deleteInstance,
  checkHealth,
} from "@/api/instances";

export const useInstanceStore = defineStore("instances", () => {
  const instances = ref([]);
  const currentInstance = ref(null);
  const loading = ref(false);
  const error = ref("");
  const health = ref({ ok: false, message: "unknown" });

  // Computed KPIs
  const totalInstances = computed(() => instances.value.length);
  const activeInstances = computed(
    () => instances.value.filter((i) => i.status === "ready").length
  );
  const failedInstances = computed(
    () => instances.value.filter((i) => i.status === "error").length
  );

  async function fetchInstances() {
    loading.value = true;
    error.value = "";
    try {
      instances.value = await listInstances();
    } catch (e) {
      error.value = e.response?.data?.error || e.message || "Failed to load";
    } finally {
      loading.value = false;
    }
  }

  async function fetchInstance(id) {
    loading.value = true;
    error.value = "";
    try {
      currentInstance.value = await getInstance(id);
    } catch (e) {
      error.value = e.response?.data?.error || e.message || "Failed to load";
    } finally {
      loading.value = false;
    }
  }

  async function addInstance(id, numInstances = 1, storageGi = 10) {
    error.value = "";
    try {
      await createInstance(id, numInstances, storageGi);
      await fetchInstances();
      return true;
    } catch (e) {
      error.value = e.response?.data?.error || e.message || "Create failed";
      return false;
    }
  }

  async function removeInstance(id) {
    error.value = "";
    try {
      await deleteInstance(id);
      await fetchInstances();
      return true;
    } catch (e) {
      error.value = e.response?.data?.error || e.message || "Delete failed";
      return false;
    }
  }

  async function fetchHealth() {
    health.value = await checkHealth();
  }

  function clearError() {
    error.value = "";
  }

  return {
    instances,
    currentInstance,
    loading,
    error,
    health,
    totalInstances,
    activeInstances,
    failedInstances,
    fetchInstances,
    fetchInstance,
    addInstance,
    removeInstance,
    fetchHealth,
    clearError,
  };
});
