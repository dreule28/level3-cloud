<!--
  InstanceDetailView — Connection info + status + copy buttons
-->
<script setup>
import { onMounted, ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useInstanceStore } from "@/stores/instances";
import GlassCard from "@/components/ui/GlassCard.vue";
import StatusPill from "@/components/ui/StatusPill.vue";
import NeonButton from "@/components/ui/NeonButton.vue";
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";

const store = useInstanceStore();
const route = useRoute();
const router = useRouter();
const showPassword = ref(false);
const copied = ref("");

const id = computed(() => route.params.id);

onMounted(() => store.fetchInstance(id.value));

function copy(text, label) {
  navigator.clipboard.writeText(text);
  copied.value = label;
  setTimeout(() => (copied.value = ""), 2000);
}

const connFields = computed(() => {
  const c = store.currentInstance?.connection;
  if (!c) return [];
  return [
    { label: "Host", value: c.host, mono: true },
    { label: "Port", value: String(c.port), mono: true },
    { label: "Database", value: c.database, mono: true },
    { label: "User", value: c.user, mono: true },
    { label: "Endpoint", value: c.endpoint, mono: true },
  ];
});

// Timeline visualization based on status
const timeline = computed(() => {
  const status = store.currentInstance?.status?.toLowerCase();
  const steps = [
    { label: "Requested", done: true },
    { label: "Provisioning", done: status === "creating" || status === "ready" },
    { label: "Ready", done: status === "ready" },
  ];
  return steps;
});
</script>

<template>
  <div class="space-y-6">
    <!-- Breadcrumb -->
    <nav class="flex items-center gap-2 text-sm text-gray-500 animate-fade-in">
      <button @click="router.push('/instances')" class="transition hover:text-neon-cyan">Instances</button>
      <span class="text-gray-700">/</span>
      <span class="font-mono text-white">{{ id }}</span>
    </nav>

    <!-- Loading -->
    <div v-if="store.loading" class="space-y-4">
      <SkeletonLoader :lines="4" height="h-6" />
    </div>

    <!-- Error -->
    <div v-else-if="store.error" class="rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-3 text-sm text-neon-red">
      {{ store.error }}
    </div>

    <!-- Content -->
    <template v-else-if="store.currentInstance">
      <!-- Header Card -->
      <GlassCard glow="glow-blue" class="p-6 animate-slide-up">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-white font-mono">{{ store.currentInstance.id }}</h1>
            <div class="mt-2">
              <StatusPill :status="store.currentInstance.status" />
            </div>
          </div>
          <NeonButton variant="ghost" @click="store.fetchInstance(id)">Refresh</NeonButton>
        </div>
      </GlassCard>

      <!-- Provisioning Timeline -->
      <GlassCard class="p-6 animate-slide-up" style="animation-delay: 100ms">
        <h3 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">Provisioning Timeline</h3>
        <div class="flex items-center gap-0">
          <template v-for="(step, i) in timeline" :key="step.label">
            <div class="flex flex-col items-center">
              <div
                :class="[
                  'flex h-10 w-10 items-center justify-center rounded-full border-2 text-xs font-bold transition-all duration-500',
                  step.done
                    ? 'border-neon-green bg-neon-green/20 text-neon-green glow-green'
                    : 'border-border text-gray-600',
                ]"
              >
                {{ step.done ? '✓' : i + 1 }}
              </div>
              <span :class="['mt-2 text-xs', step.done ? 'text-neon-green' : 'text-gray-600']">
                {{ step.label }}
              </span>
            </div>
            <div
              v-if="i < timeline.length - 1"
              :class="[
                'mb-5 h-0.5 flex-1 transition-all duration-500',
                step.done ? 'bg-gradient-to-r from-neon-green to-neon-green/30' : 'bg-border',
              ]"
            />
          </template>
        </div>
      </GlassCard>

      <!-- Connection Info -->
      <GlassCard
        v-if="store.currentInstance.connection"
        glow="glow-cyan"
        class="p-6 animate-slide-up"
        style="animation-delay: 200ms"
      >
        <h3 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">Connection Details</h3>

        <div class="space-y-3">
          <div
            v-for="field in connFields"
            :key="field.label"
            class="flex items-center justify-between rounded-lg border border-border/50 bg-void/30 px-4 py-3"
          >
            <div>
              <span class="text-xs text-gray-500">{{ field.label }}</span>
              <p :class="['text-sm text-white', field.mono && 'font-mono']">{{ field.value }}</p>
            </div>
            <button
              @click="copy(field.value, field.label)"
              :class="[
                'rounded-md border border-border px-3 py-1 text-xs transition',
                copied === field.label
                  ? 'border-neon-green/50 text-neon-green'
                  : 'text-gray-500 hover:border-neon-cyan/50 hover:text-neon-cyan',
              ]"
            >
              {{ copied === field.label ? '✓ Copied' : 'Copy' }}
            </button>
          </div>

          <!-- Password (hidden by default) -->
          <div class="flex items-center justify-between rounded-lg border border-border/50 bg-void/30 px-4 py-3">
            <div>
              <span class="text-xs text-gray-500">Password</span>
              <p class="text-sm font-mono text-white">
                {{ showPassword ? store.currentInstance.connection.password : '••••••••••••' }}
              </p>
            </div>
            <div class="flex gap-2">
              <button
                @click="showPassword = !showPassword"
                class="rounded-md border border-border px-3 py-1 text-xs text-gray-500 transition hover:border-neon-purple/50 hover:text-neon-purple"
              >
                {{ showPassword ? 'Hide' : 'Show' }}
              </button>
              <button
                @click="copy(store.currentInstance.connection.password, 'Password')"
                :class="[
                  'rounded-md border border-border px-3 py-1 text-xs transition',
                  copied === 'Password'
                    ? 'border-neon-green/50 text-neon-green'
                    : 'text-gray-500 hover:border-neon-cyan/50 hover:text-neon-cyan',
                ]"
              >
                {{ copied === 'Password' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
          </div>
        </div>
      </GlassCard>

      <!-- No connection yet -->
      <GlassCard
        v-else
        class="flex flex-col items-center justify-center py-12 animate-slide-up"
        style="animation-delay: 200ms"
      >
        <span class="mb-3 text-3xl animate-glow-pulse">⏳</span>
        <p class="text-sm text-gray-500">Connection details will appear once the instance is ready.</p>
      </GlassCard>
    </template>
  </div>
</template>
