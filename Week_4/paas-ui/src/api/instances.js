/**
 * Instances API — CRUD + health check
 */
import api from "./axios";

export async function listInstances() {
  const { data } = await api.get("/instances");
  return data;
}

export async function getInstance(id) {
  const { data } = await api.get(`/instances/${encodeURIComponent(id)}`);
  return data;
}

export async function createInstance(id, instances = 1, storageGi = 10) {
  const { data } = await api.post("/instances", { id, instances, storageGi });
  return data;
}

export async function deleteInstance(id) {
  await api.delete(`/instances/${encodeURIComponent(id)}`);
}

export async function checkHealth() {
  try {
    const { data } = await api.get("/healthz");
    return { ok: true, message: data };
  } catch {
    return { ok: false, message: "unreachable" };
  }
}
