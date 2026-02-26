<!--
  InstanceDetailView — GSAP-animated connection info + timeline + copy-to-clipboard + toast
-->
<script setup>
import { computed, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useInstanceStore } from "@/stores/instances";
import { getInstanceLogs } from "@/api/logs";
import { useToast } from "@/composables/useToast";
import GlassCard from "@/components/ui/GlassCard.vue";
import StatusPill from "@/components/ui/StatusPill.vue";
import NeonButton from "@/components/ui/NeonButton.vue";
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";

const store = useInstanceStore();
const route = useRoute();
const router = useRouter();
const toast = useToast();
const showPassword = ref(false);
const copied = ref("");
const logs = ref([]);
const logsLoading = ref(false);
const logsError = ref("");
const logType = ref("all");
const POLL_INTERVAL_MS = 3000;
const LOG_LIMIT = 100;
const logTypeOptions = [
  { label: "All", value: "all" },
  { label: "Audit", value: "audit" },
  { label: "Service", value: "service" },
];

const id = computed(() => route.params.id);
const normalizedStatus = computed(
  () => store.currentInstance?.status?.toLowerCase() || ""
);
const shouldAutoRefresh = computed(() => {
  const inst = store.currentInstance;
  if (!inst || inst.id !== id.value) return false;

  const status = (inst.status || "").toLowerCase();
  if (status === "error") return false;
  if (status === "ready") return !inst.connection;

  return ["creating", "pending", "provisioning", "requested"].includes(status);
});
let pollTimer = null;
let pollInFlight = false;

watch(
  id,
  async (nextId, prevId) => {
    if (!nextId) return;
    if (nextId !== prevId) {
      showPassword.value = false;
      copied.value = "";
      logs.value = [];
      logsError.value = "";
      stopPolling();
    }
    await Promise.allSettled([store.fetchInstance(nextId), fetchLogs({ silent: false })]);
  },
  { immediate: true }
);

watch(logType, () => {
  if (!id.value) return;
  void fetchLogs();
});

watch(
  shouldAutoRefresh,
  (enabled) => {
    if (enabled) {
      startPolling();
      return;
    }
    stopPolling();
  },
  { immediate: true }
);

onUnmounted(stopPolling);

function startPolling() {
  if (pollTimer) return;
  void pollInstance();
  pollTimer = window.setInterval(() => {
    void pollInstance();
  }, POLL_INTERVAL_MS);
}

function stopPolling() {
  if (!pollTimer) return;
  clearInterval(pollTimer);
  pollTimer = null;
  pollInFlight = false;
}

async function pollInstance() {
  if (pollInFlight || !id.value) return;
  pollInFlight = true;
  try {
    await Promise.allSettled([store.fetchInstance(id.value), fetchLogs({ silent: true })]);
  } finally {
    pollInFlight = false;
  }
}

async function fetchLogs({ silent = false } = {}) {
  if (!id.value) return;

  if (!silent) {
    logsLoading.value = true;
  }
  logsError.value = "";

  try {
    const payload = await getInstanceLogs(id.value, {
      type: logType.value === "all" ? "" : logType.value,
      limit: LOG_LIMIT,
    });
    logs.value = normalizeLogsPayload(payload);
  } catch (e) {
    logsError.value = e.response?.data?.error || e.message || "Failed to load logs";
  } finally {
    if (!silent) {
      logsLoading.value = false;
    }
  }
}

function normalizeLogsPayload(payload) {
  if (payload == null) {
    return [];
  }

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
      const looksHtml = payload.includes("<html") || payload.includes("<!doctype html");
      throw new Error(
        looksHtml
          ? "Logs endpoint returned HTML (likely old backend/proxy response)"
          : "Unexpected logs response format (string)"
      );
    }

    if (payload && typeof payload === "object") {
      if (payload.error) {
        throw new Error(String(payload.error));
      }
      if (payload.message) {
        throw new Error(String(payload.message));
      }

      // Last-resort support for object maps with numeric keys.
      const values = Object.values(payload);
      if (values.length && values.every((v) => v && typeof v === "object")) {
        return values.map((raw, i) => normalizeLogEntry(raw, i));
      }

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

function legacyCopy(text) {
  const textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "");
  textarea.style.position = "fixed";
  textarea.style.opacity = "0";
  textarea.style.pointerEvents = "none";
  document.body.appendChild(textarea);
  textarea.select();
  textarea.setSelectionRange(0, textarea.value.length);

  let copiedOk = false;
  try {
    copiedOk = document.execCommand("copy");
  } catch {
    copiedOk = false;
  } finally {
    document.body.removeChild(textarea);
  }

  return copiedOk;
}

async function copy(text, label) {
  const value = String(text ?? "");
  if (!value) {
    toast.error(`${label} is empty`);
    return;
  }

  try {
    if (window.isSecureContext && navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(value);
    } else if (!legacyCopy(value)) {
      throw new Error("Clipboard API unavailable");
    }

    copied.value = label;
    toast.success(`${label} copied to clipboard`);
    setTimeout(() => (copied.value = ""), 2000);
  } catch {
    toast.error(`Couldn't copy ${label}. Select the text manually.`);
  }
}

const connFields = computed(() => {
  const c = store.currentInstance?.connection;
  if (!c) return [];
  return [
    { label: "Host", value: c.host, mono: true, icon: "🖥" },
    { label: "Port", value: String(c.port), mono: true, icon: "🔌" },
    { label: "Database", value: c.database, mono: true, icon: "🗄" },
    { label: "User", value: c.user, mono: true, icon: "👤" },
    { label: "Endpoint", value: c.endpoint, mono: true, icon: "🌐" },
  ];
});

const timeline = computed(() => {
  const status = store.currentInstance?.status?.toLowerCase();
  return [
    { label: "Requested", done: true, icon: "📋" },
    { label: "Provisioning", done: status === "creating" || status === "ready", icon: "⚙" },
    { label: "Ready", done: status === "ready", icon: "✓" },
  ];
});

// psql connection string
const psqlString = computed(() => {
  const c = store.currentInstance?.connection;
  if (!c) return "";
  return `psql -h ${c.host} -p ${c.port} -U ${c.user} -d ${c.database}`;
});

async function refreshNow() {
  await Promise.allSettled([store.fetchInstance(id.value), fetchLogs({ silent: false })]);
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
    <!-- Breadcrumb -->
    <nav class="flex items-center gap-2 text-sm text-gray-500">
      <button @click="router.push('/instances')" class="transition hover:text-neon-cyan">
        ← Instances
      </button>
      <span class="text-gray-700">/</span>
      <span class="font-mono text-neon-cyan">{{ id }}</span>
    </nav>

    <!-- Loading -->
    <div v-if="store.loading && !store.currentInstance" class="space-y-4">
      <SkeletonLoader :lines="4" height="h-6" />
    </div>

    <!-- Error -->
    <div v-else-if="store.error && !store.currentInstance" class="detail-card rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-3 text-sm text-neon-red">
      {{ store.error }}
    </div>

    <!-- Content -->
    <template v-else-if="store.currentInstance">
      <!-- Header Card -->
      <GlassCard glow="glow-blue" class="detail-card p-6">
        <div class="flex items-center justify-between">
          <div>
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-neon-blue to-neon-purple text-lg shadow-[0_0_20px_rgba(59,130,246,0.2)]">
                ⬡
              </div>
              <div>
                <h1 class="text-2xl font-bold text-white font-mono">{{ store.currentInstance.id }}</h1>
                <p class="text-xs font-mono text-gray-600">cluster.cloudnativepg.io</p>
              </div>
            </div>
            <div class="mt-3">
              <StatusPill :status="store.currentInstance.status" />
              <p
                v-if="shouldAutoRefresh"
                class="mt-2 text-[11px] font-mono text-gray-500"
              >
                Auto-refreshing every 3s while provisioning...
              </p>
            </div>
          </div>
          <NeonButton variant="ghost" @click="refreshNow">↻ Refresh</NeonButton>
        </div>
      </GlassCard>

      <!-- Provisioning Timeline -->
      <GlassCard class="detail-card p-6">
        <h3 class="mb-5 text-xs font-semibold uppercase tracking-wider text-gray-500">Provisioning Timeline</h3>
        <div class="flex items-center gap-0">
          <template v-for="(step, i) in timeline" :key="step.label">
            <div class="flex flex-col items-center">
              <div
                :class="[
                  'flex h-12 w-12 items-center justify-center rounded-full border-2 text-sm font-bold transition-all duration-500',
                  step.done
                    ? 'border-neon-green bg-neon-green/15 text-neon-green shadow-[0_0_15px_rgba(16,185,129,0.2)]'
                    : 'border-border text-gray-600',
                ]"
              >
                {{ step.done ? step.icon : i + 1 }}
              </div>
              <span :class="['mt-2 text-xs font-medium', step.done ? 'text-neon-green' : 'text-gray-600']">
                {{ step.label }}
              </span>
            </div>
            <div
              v-if="i < timeline.length - 1"
              :class="[
                'mb-6 h-0.5 flex-1 transition-all duration-700',
                step.done ? 'bg-gradient-to-r from-neon-green to-neon-green/30 shadow-[0_0_6px_rgba(16,185,129,0.3)]' : 'bg-border',
              ]"
            />
          </template>
        </div>
      </GlassCard>

      <!-- Connection Info -->
      <GlassCard
        v-if="store.currentInstance.connection"
        glow="glow-cyan"
        class="detail-card p-6"
      >
        <h3 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">Connection Details</h3>

        <div class="space-y-2">
          <div
            v-for="field in connFields"
            :key="field.label"
            class="group flex items-center justify-between rounded-lg border border-border/50 bg-void/30 px-4 py-3 transition-all duration-200 hover:border-neon-cyan/20 hover:bg-void/50"
          >
            <div class="flex items-center gap-3">
              <span class="text-base">{{ field.icon }}</span>
              <div>
                <span class="text-[10px] uppercase tracking-wider text-gray-600">{{ field.label }}</span>
                <p :class="['select-text text-sm text-white', field.mono && 'font-mono']">{{ field.value }}</p>
              </div>
            </div>
            <button
              @click="copy(field.value, field.label)"
              type="button"
              :class="[
                'rounded-md border px-3 py-1 text-xs font-medium transition-all duration-200',
                copied === field.label
                  ? 'border-neon-green/50 text-neon-green bg-neon-green/5'
                  : 'border-border text-gray-500 hover:border-neon-cyan/50 hover:text-neon-cyan',
              ]"
            >
              {{ copied === field.label ? '✓ Copied' : 'Copy' }}
            </button>
          </div>

          <!-- Password (hidden by default) -->
          <div class="group flex items-center justify-between rounded-lg border border-border/50 bg-void/30 px-4 py-3 transition-all duration-200 hover:border-neon-purple/20">
            <div class="flex items-center gap-3">
              <span class="text-base">🔑</span>
              <div>
                <span class="text-[10px] uppercase tracking-wider text-gray-600">Password</span>
                <p class="select-text text-sm font-mono text-white">
                  {{ showPassword ? store.currentInstance.connection.password : '••••••••••••' }}
                </p>
              </div>
            </div>
            <div class="flex gap-2">
              <button
                @click="showPassword = !showPassword"
                type="button"
                class="rounded-md border border-border px-3 py-1 text-xs font-medium text-gray-500 transition hover:border-neon-purple/50 hover:text-neon-purple"
              >
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
              <button
                @click="copy(store.currentInstance.connection.password, 'Password')"
                type="button"
                :class="[
                  'rounded-md border px-3 py-1 text-xs font-medium transition-all duration-200',
                  copied === 'Password'
                    ? 'border-neon-green/50 text-neon-green bg-neon-green/5'
                    : 'border-border text-gray-500 hover:border-neon-cyan/50 hover:text-neon-cyan',
                ]"
              >
                {{ copied === 'Password' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Quick connect command -->
        <div class="mt-4 rounded-lg border border-border/50 bg-void/60 px-4 py-3">
          <div class="flex items-center justify-between mb-1">
            <span class="text-[10px] uppercase tracking-wider text-gray-600">Quick Connect</span>
            <button
              @click="copy(psqlString, 'psql')"
              type="button"
              :class="[
                'text-[10px] font-medium transition',
                copied === 'psql' ? 'text-neon-green' : 'text-gray-600 hover:text-neon-cyan',
              ]"
            >
              {{ copied === 'psql' ? '✓ Copied' : 'Copy' }}
            </button>
          </div>
          <code class="block select-text break-all font-mono text-xs text-neon-cyan">{{ psqlString }}</code>
        </div>
      </GlassCard>

      <!-- No connection yet -->
      <GlassCard
        v-else
        class="detail-card flex flex-col items-center justify-center py-16"
      >
        <span class="mb-3 text-4xl animate-glow-pulse">⏳</span>
        <p class="text-sm text-gray-500 font-medium">Provisioning in progress</p>
        <p class="mt-1 text-xs text-gray-600">Connection details appear once the instance is ready. This page auto-refreshes every 3s.</p>
        <NeonButton variant="ghost" class="mt-4" @click="refreshNow">
          ↻ Check Status
        </NeonButton>
      </GlassCard>

      <!-- User-facing logs -->
      <GlassCard class="detail-card p-6">
        <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h3 class="text-xs font-semibold uppercase tracking-wider text-gray-500">Instance Logs</h3>
            <p class="mt-1 text-xs text-gray-600">
              Audit = user-triggered actions. Service = platform/system events (accepted requests, status transitions).
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button
              v-for="option in logTypeOptions"
              :key="option.value"
              type="button"
              @click="logType = option.value"
              :class="[
                'rounded-md border px-3 py-1 text-xs font-medium transition',
                logType === option.value
                  ? 'border-neon-cyan/60 bg-neon-cyan/10 text-neon-cyan'
                  : 'border-border text-gray-500 hover:border-neon-cyan/30 hover:text-white',
              ]"
            >
              {{ option.label }}
            </button>
            <NeonButton variant="ghost" @click="fetchLogs()" :loading="logsLoading && logs.length > 0">
              ↻ Refresh Logs
            </NeonButton>
          </div>
        </div>

        <div v-if="logsError" class="rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-3 text-xs text-neon-red">
          {{ logsError }}
        </div>

        <div v-else-if="logsLoading && logs.length === 0" class="space-y-3">
          <SkeletonLoader :lines="4" height="h-5" />
        </div>

        <div
          v-else-if="logs.length === 0"
          class="rounded-lg border border-border/40 bg-void/30 px-4 py-6 text-center text-sm text-gray-500"
        >
          No logs yet for this instance.
        </div>

        <div v-else class="max-h-[26rem] space-y-2 overflow-y-auto pr-1">
          <div
            v-for="entry in logs"
            :key="entry.id"
            class="rounded-lg border border-border/40 bg-void/35 px-4 py-3"
          >
            <div class="flex flex-wrap items-center gap-2">
              <span
                :class="[
                  'rounded-full border px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider',
                  entry.type === 'audit'
                    ? 'border-neon-purple/40 bg-neon-purple/10 text-neon-purple'
                    : 'border-neon-blue/40 bg-neon-blue/10 text-neon-blue',
                ]"
              >
                {{ entry.type === 'audit' ? 'User' : 'System' }}
              </span>
              <span class="font-mono text-[11px] text-gray-500">{{ formatLogTime(entry.timestamp) }}</span>
              <span class="font-mono text-[11px] text-neon-cyan">{{ entry.action }}</span>
              <span v-if="entry.user" class="font-mono text-[11px] text-gray-600">user={{ entry.user }}</span>
            </div>
            <p class="mt-1 select-text text-sm text-gray-300">{{ entry.message }}</p>
          </div>
        </div>
      </GlassCard>
    </template>
  </div>
</template>
