<!--
  InstanceDetailView — GSAP-animated connection info + timeline + copy-to-clipboard + toast
-->
<script setup>
import { onMounted, ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useInstanceStore } from "@/stores/instances";
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

const id = computed(() => route.params.id);

onMounted(() => store.fetchInstance(id.value));

function copy(text, label) {
  navigator.clipboard.writeText(text);
  copied.value = label;
  toast.success(`${label} copied to clipboard`);
  setTimeout(() => (copied.value = ""), 2000);
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
    <div v-if="store.loading" class="space-y-4">
      <SkeletonLoader :lines="4" height="h-6" />
    </div>

    <!-- Error -->
    <div v-else-if="store.error" class="detail-card rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-3 text-sm text-neon-red">
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
            </div>
          </div>
          <NeonButton variant="ghost" @click="store.fetchInstance(id)">↻ Refresh</NeonButton>
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
                <p :class="['text-sm text-white', field.mono && 'font-mono']">{{ field.value }}</p>
              </div>
            </div>
            <button
              @click="copy(field.value, field.label)"
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
                <p class="text-sm font-mono text-white">
                  {{ showPassword ? store.currentInstance.connection.password : '••••••••••••' }}
                </p>
              </div>
            </div>
            <div class="flex gap-2">
              <button
                @click="showPassword = !showPassword"
                class="rounded-md border border-border px-3 py-1 text-xs font-medium text-gray-500 transition hover:border-neon-purple/50 hover:text-neon-purple"
              >
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
              <button
                @click="copy(store.currentInstance.connection.password, 'Password')"
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
              :class="[
                'text-[10px] font-medium transition',
                copied === 'psql' ? 'text-neon-green' : 'text-gray-600 hover:text-neon-cyan',
              ]"
            >
              {{ copied === 'psql' ? '✓ Copied' : 'Copy' }}
            </button>
          </div>
          <code class="block font-mono text-xs text-neon-cyan break-all">{{ psqlString }}</code>
        </div>
      </GlassCard>

      <!-- No connection yet -->
      <GlassCard
        v-else
        class="detail-card flex flex-col items-center justify-center py-16"
      >
        <span class="mb-3 text-4xl animate-glow-pulse">⏳</span>
        <p class="text-sm text-gray-500 font-medium">Provisioning in progress</p>
        <p class="mt-1 text-xs text-gray-600">Connection details appear once the instance is ready.</p>
        <NeonButton variant="ghost" class="mt-4" @click="store.fetchInstance(id)">
          ↻ Check Status
        </NeonButton>
      </GlassCard>
    </template>
  </div>
</template>
