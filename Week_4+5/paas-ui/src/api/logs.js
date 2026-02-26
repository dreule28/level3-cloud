import api from "./axios";

export async function getLogs({ type = "", instanceId = "", limit = 100 } = {}) {
  const params = {};
  if (type) params.type = type;
  if (instanceId) params.instanceId = instanceId;
  if (limit) params.limit = limit;

  const { data } = await api.get("/logs", { params });
  return data;
}

export async function getInstanceLogs(id, { type = "", limit = 100 } = {}) {
  const params = {};
  if (type) params.type = type;
  if (limit) params.limit = limit;

  const { data } = await api.get(`/instances/${encodeURIComponent(id)}/logs`, { params });
  return data;
}
