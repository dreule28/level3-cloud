<!--
  DashboardView — Animated KPI cards + cluster health
-->
<script setup>
import { onMounted } from "vue";
import { useInstanceStore } from "@/stores/instances";
import GlassCard from "@/components/ui/GlassCard.vue";
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";

const store = useInstanceStore();

onMounted(async () => {
  await Promise.all([store.fetchInstances(), store.fetchHealth()]);
});

const kpis = [
  { key: "active", label: "Active", icon: "⚡", color: "neon-green", glow: "glow-green" },
  { key: "total", label: "Total", icon: "⬡", color: "neon-blue", glow: "glow-blue" },
  { key: "failed", label: "Failed", icon: "⚠", color: "neon-red", glow: "glow-red" },
];

function getKpiValue(key) {
  if (key === "active") return store.activeInstances;
  if (key === "total") return store.totalInstances;
  if (key === "failed") return store.failedInstances;
  return 0;
}
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="animate-fade-in">
      <h1 class="text-3xl font-bold text-white">Dashboard</h1>
      <p class="mt-1 text-sm text-gray-500">PaaS control plane overview</p>
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <GlassCard
        v-for="(kpi, i) in kpis"
        :key="kpi.key"
        :glow="kpi.glow"
        class="p-6"
        :style="{ animationDelay: `${i * 100}ms` }"
      >
        <div v-if="store.loading" class="space-y-3">
          <SkeletonLoader :lines="2" />
        </div>
        <div v-else class="animate-slide-up" :style="{ animationDelay: `${i * 100}ms` }">
          <div class="flex items-center justify-between">
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">{{ kpi.label }}</span>
            <span class="text-lg">{{ kpi.icon }}</span>
          </div>
          <p :class="['mt-2 text-4xl font-bold', `text-${kpi.color}`]">
            {{ getKpiValue(kpi.key) }}
          </p>
        </div>
      </GlassCard>

      <!-- API Health -->
      <GlassCard :glow="store.health.ok ? 'glow-green' : 'glow-red'" class="p-6">
        <div v-if="store.loading" class="space-y-3">
          <SkeletonLoader :lines="2" />
        </div>
        <div v-else class="animate-slide-up" style="animation-delay: 300ms">
          <div class="flex items-center justify-between">
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">API Health</span>
            <span class="text-lg">♥</span>
          </div>
          <div class="mt-2 flex items-center gap-2">
            <span
              :class="[
                'h-3 w-3 rounded-full',
                store.health.ok ? 'bg-neon-green animate-glow-pulse' : 'bg-neon-red animate-flicker',
              ]"
            />
            <span :class="['text-lg font-bold', store.health.ok ? 'text-neon-green' : 'text-neon-red']">
              {{ store.health.ok ? "Healthy" : "Down" }}
            </span>
          </div>
        </div>
      </GlassCard>
    </div>

    <!-- Quick links -->
    <GlassCard class="p-6 animate-slide-up" style="animation-delay: 400ms">
      <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-4">Quick Actions</h3>
      <div class="flex flex-wrap gap-3">
        <router-link
          to="/instances"
          class="rounded-lg border border-border bg-surface/50 px-4 py-2 text-sm text-gray-300 transition hover:border-neon-blue/50 hover:text-white hover:glow-blue"
        >
          ⬡ View Instances
        </router-link>
      </div>
    </GlassCard>
  </div>
</template>
