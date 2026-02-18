import { http } from "./http";

export function listInstances() {
  return http.get("/instances");
}

export function getInstance(id) {
  return http.get(`/instances/${encodeURIComponent(id)}`);
}

export function createInstance(id) {
  return http.post("/instances", { id });
}

export function deleteInstance(id) {
  return http.del(`/instances/${encodeURIComponent(id)}`);
}
