import { http } from "./http";

export function listInstances() {
  return http.get("/api/instances");
}

export function getInstance(id) {
  return http.get(`/api/instances/${encodeURIComponent(id)}`);
}

export function createInstance(id) {
  return http.post("/api/instances", { id });
}

export function deleteInstance(id) {
  return http.del(`/api/instances/${encodeURIComponent(id)}`);
}
