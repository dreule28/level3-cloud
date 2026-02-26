<!--
  LogsView — Global audit/service log stream across all instances
-->
<script setup>
import { computed, onMounted, ref } from "vue";
import { getLogs } from "@/api/logs";
import { useToast } from "@/composables/useToast";
import GlassCard from "@/components/ui/GlassCard.vue";
import NeonButton from "@/components/ui/NeonButton.vue";
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";

const toast = useToast();
const logs = ref([]);
const loading = ref(false);
const error = ref("");
const typeFilter = ref("all");
const opFilter = ref("all");
const instanceFilter = ref("");

const typeOptions = [
  { label: "All", value: "all" },
  { label: "Audit", value: "audit" },
  { label: "Service", value: "service" },
];

const opOptions = [
  { label: "All Ops", value: "all" },
  { label: "Create", value: "create" },
  { label: "Delete", value: "delete" },
  { label: "Status", value: "status" },
  { label: "Other", value: "other" },
];

const visibleLogs = computed(() => {
  const instanceNeedle = instanceFilter.value.trim().toLowerCase();

  return logs.value.filter((entry) => {
    if (instanceNeedle && !String(entry.instanceId || "").toLowerCase().includes(instanceNeedle)) {
      return false;
    }

    if (opFilter.value !== "all" && operationKind(entry) !== opFilter.value) {
      return false;
    }

    return true;
  });
});

onMounted(() => {
  void fetchAllLogs();
});

async function fetchAllLogs() {
  loading.value = true;
  error.value = "";
  try {
    const payload = await getLogs({
      type: typeFilter.value === "all" ? "" : typeFilter.value,
      limit: 200,
    });
    logs.value = normalizeLogsPayload(payload);
  } catch (e) {
    error.value = e.response?.data?.error || e.message || "Failed to load logs";
    toast.error(error.value);
  } finally {
    loading.value = false;
  }
}

function operationKind(entry) {
  const action = String(entry?.action || "").toLowerCase();
  if (action.includes(".create.")) return "create";
  if (action.includes(".delete.")) return "delete";
  if (action.includes("status.")) return "status";
  return "other";
}

function operationLabel(entry) {
  const kind = operationKind(entry);
  if (kind === "create") return "CREATE";
  if (kind === "delete") return "DELETE";
  if (kind === "status") return "STATUS";
  return "OTHER";
}

function operationBadgeClass(entry) {
  const kind = operationKind(entry);
  if (kind === "create") {
    return "border-neon-green/40 bg-neon-green/10 text-neon-green";
  }
  if (kind === "delete") {
    return "border-neon-red/40 bg-neon-red/10 text-neon-red";
  }
  if (kind === "status") {
    return "border-neon-amber/40 bg-neon-amber/10 text-neon-amber";
  }
  return "border-border bg-void/40 text-gray-300";
}

function actionCodeClass(entry) {
  const kind = operationKind(entry);
  if (kind === "create") return "text-neon-green/80";
  if (kind === "delete") return "text-neon-red/80";
  if (kind === "status") return "text-neon-amber/80";
  return "text-gray-600";
}

function normalizeLogsPayload(payload) {
  if (payload == null) return [];

  const list = Array.isArray(payload)
    ? payload
    : Array.isArray(payload?.items)
      ? payload.items
      : Array.isArray(payload?.logs)
        ? payload.logs
        : Array.isArray(payload?.entries)
          ? payload.entries
          : Array.isArray(payload?.data)
            ? payload.data
            : null;

  if (!list) {
    if (typeof payload === "string") {
      throw new Error(
        payload.includes("<html") || payload.includes("<!doctype html")
          ? "Logs endpoint returned HTML (wrong backend/proxy target)"
          : "Unexpected logs response format (string)"
      );
    }
    if (payload && typeof payload === "object") {
      if (payload.error) throw new Error(String(payload.error));
      if (payload.message) throw new Error(String(payload.message));
      throw new Error(`Unexpected logs response format (object keys: ${Object.keys(payload).slice(0, 5).join(", ") || "none"})`);
    }
    throw new Error(`Unexpected logs response format (${typeof payload})`);
  }

  return list
    .map((raw, i) => normalizeLogEntry(raw, i))
    .filter((entry) => entry.type || entry.action || entry.message || entry.timestamp);
}

function normalizeLogEntry(raw, index) {
  if (!raw || typeof raw !== "object") {
    return {
      id: `raw-${index}`,
      instanceId: "",
      type: "",
      action: "",
      message: String(raw ?? ""),
      user: "",
      timestamp: "",
    };
  }

  return {
    id: raw.id ?? raw.ID ?? `log-${index}`,
    instanceId: raw.instanceId ?? raw.instance_id ?? raw.InstanceID ?? "",
    type: raw.type ?? raw.Type ?? "",
    action: raw.action ?? raw.Action ?? "",
    message: raw.message ?? raw.Message ?? "",
    user: raw.user ?? raw.User ?? "",
    timestamp: raw.timestamp ?? raw.Timestamp ?? "",
  };
}

function formatLogTime(value) {
  if (!value) return "unknown";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString();
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-white">Logs</h1>
        <p class="mt-1 text-sm font-mono text-gray-500">
          <span class="text-neon-cyan">audit</span>/<span class="text-neon-blue">service</span> — global event stream
        </p>
      </div>
      <NeonButton variant="ghost" @click="fetchAllLogs" :loading="loading">
        ↻ Refresh Logs
      </NeonButton>
    </div>

    <GlassCard class="p-5">
      <div class="grid gap-4 lg:grid-cols-[1fr_auto_auto]">
        <label class="block">
          <span class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Filter by Instance</span>
          <input
            v-model="instanceFilter"
            type="text"
            placeholder="e.g. pg-demo"
            class="w-full rounded-lg border border-border bg-void/50 px-4 py-2.5 font-mono text-sm text-white placeholder-gray-600 transition focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50"
          />
        </label>

        <div>
          <span class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Log Type</span>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="option in typeOptions"
              :key="option.value"
              type="button"
              @click="typeFilter = option.value; fetchAllLogs()"
              :class="[
                'rounded-md border px-3 py-2 text-xs font-medium transition',
                typeFilter === option.value
                  ? 'border-neon-cyan/60 bg-neon-cyan/10 text-neon-cyan'
                  : 'border-border text-gray-500 hover:border-neon-cyan/30 hover:text-white',
              ]"
            >
              {{ option.label }}
            </button>
          </div>
        </div>

        <div>
          <span class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Operation</span>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="option in opOptions"
              :key="option.value"
              type="button"
              @click="opFilter = option.value"
              :class="[
                'rounded-md border px-3 py-2 text-xs font-medium transition',
                opFilter === option.value
                  ? 'border-neon-blue/60 bg-neon-blue/10 text-neon-blue'
                  : 'border-border text-gray-500 hover:border-neon-blue/30 hover:text-white',
              ]"
            >
              {{ option.label }}
            </button>
          </div>
        </div>
      </div>

      <div class="mt-4 rounded-lg border border-border/40 bg-void/25 px-4 py-3 text-xs text-gray-400">
        <span class="text-neon-purple font-semibold">Audit</span> = user-triggered actions (who requested create/delete).
        <span class="mx-2 text-gray-600">|</span>
        <span class="text-neon-blue font-semibold">Service</span> = platform/system events (accepted request, observed status changes).
      </div>
    </GlassCard>

    <GlassCard class="overflow-hidden">
      <div class="grid grid-cols-12 gap-4 border-b border-border px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-600">
        <div class="col-span-2">Type</div>
        <div class="col-span-2">Operation</div>
        <div class="col-span-2">Instance</div>
        <div class="col-span-2">User</div>
        <div class="col-span-4">Message / Time</div>
      </div>

      <div v-if="error" class="px-6 py-4 text-sm text-neon-red">
        {{ error }}
      </div>

      <div v-else-if="loading && logs.length === 0" class="p-6 space-y-4">
        <SkeletonLoader :lines="5" height="h-10" />
      </div>

      <div v-else-if="visibleLogs.length === 0" class="px-6 py-12 text-center text-sm text-gray-500">
        No logs match the current filters.
      </div>

      <div v-else class="max-h-[34rem] overflow-y-auto">
        <div
          v-for="entry in visibleLogs"
          :key="entry.id"
          class="grid grid-cols-12 gap-4 border-b border-border/20 px-6 py-4 transition hover:bg-white/[0.02]"
        >
          <div class="col-span-2">
            <span
              :class="[
                'inline-flex rounded-full border px-2 py-1 text-[10px] font-semibold uppercase tracking-wider',
                entry.type === 'audit'
                  ? 'border-neon-purple/40 bg-neon-purple/10 text-neon-purple'
                  : 'border-neon-blue/40 bg-neon-blue/10 text-neon-blue',
              ]"
            >
              {{ entry.type === 'audit' ? 'User' : 'System' }}
            </span>
          </div>

          <div class="col-span-2">
            <span
              :class="[
                'inline-flex rounded-md border px-2 py-1 font-mono text-[10px]',
                operationBadgeClass(entry),
              ]"
            >
              {{ operationLabel(entry) }}
            </span>
            <p :class="['mt-1 break-all font-mono text-[11px]', actionCodeClass(entry)]">{{ entry.action }}</p>
          </div>

          <div class="col-span-2">
            <span class="break-all font-mono text-xs text-neon-cyan">{{ entry.instanceId || "global" }}</span>
          </div>

          <div class="col-span-2">
            <span class="font-mono text-xs text-gray-400">{{ entry.user || "system" }}</span>
          </div>

          <div class="col-span-4">
            <p class="select-text text-sm text-gray-200">{{ entry.message || "(no message)" }}</p>
            <p class="mt-1 font-mono text-[11px] text-gray-600">{{ formatLogTime(entry.timestamp) }}</p>
          </div>
        </div>
      </div>
    </GlassCard>
  </div>
</template>
