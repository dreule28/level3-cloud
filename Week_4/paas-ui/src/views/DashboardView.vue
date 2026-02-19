<!--
  DashboardView — Animated KPI cards with GSAP count-up + cluster health HUD
-->
<script setup>
import { onMounted, ref, watch } from "vue";
import gsap from "gsap";
import { useInstanceStore } from "@/stores/instances";
import GlassCard from "@/components/ui/GlassCard.vue";
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";
import NeonButton from "@/components/ui/NeonButton.vue";

const store = useInstanceStore();

const animatedValues = ref({ active: 0, total: 0, failed: 0 });

onMounted(async () => {
  await Promise.all([store.fetchInstances(), store.fetchHealth()]);
});

// Animate numbers counting up when data changes
watch(
  () => [store.activeInstances, store.totalInstances, store.failedInstances],
  ([active, total, failed]) => {
    gsap.to(animatedValues.value, { active, duration: 1.2, ease: "power2.out" });
    gsap.to(animatedValues.value, { total, duration: 1.2, ease: "power2.out" });
    gsap.to(animatedValues.value, { failed, duration: 1.2, ease: "power2.out" });
  },
  { immediate: true }
);

const kpis = [
  { key: "active", label: "Active", icon: "⚡", glow: "glow-green", gradient: "from-neon-green/20 to-transparent", textClass: "text-neon-green", barClass: "bg-neon-green/30" },
  { key: "total", label: "Total", icon: "⬡", glow: "glow-blue", gradient: "from-neon-blue/20 to-transparent", textClass: "text-neon-blue", barClass: "bg-neon-blue/30" },
  { key: "failed", label: "Failed", icon: "⚠", glow: "glow-red", gradient: "from-neon-red/20 to-transparent", textClass: "text-neon-red", barClass: "bg-neon-red/30" },
];

function getKpiValue(key) {
  return Math.round(animatedValues.value[key] || 0);
}
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold text-white">Dashboard</h1>
      <p class="mt-1 text-sm text-gray-500 font-mono">
        <span class="text-neon-cyan">sys</span>.overview — PaaS control plane
      </p>
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <GlassCard
        v-for="kpi in kpis"
        :key="kpi.key"
        :glow="kpi.glow"
        class="kpi-card group relative overflow-hidden p-6"
      >
        <!-- Background gradient pulse -->
        <div :class="['absolute inset-0 bg-gradient-to-br opacity-0 transition-opacity duration-500 group-hover:opacity-100', kpi.gradient]" />

        <div v-if="store.loading" class="relative space-y-3">
          <SkeletonLoader :lines="2" />
        </div>
        <div v-else class="relative">
          <div class="flex items-center justify-between">
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">{{ kpi.label }}</span>
            <span class="text-lg transition-transform duration-300 group-hover:scale-125">{{ kpi.icon }}</span>
          </div>
          <p :class="['mt-2 text-5xl font-black tabular-nums', kpi.textClass]">
            {{ getKpiValue(kpi.key) }}
          </p>
          <!-- Mini sparkline decoration -->
          <div class="mt-3 flex items-end gap-0.5 h-4">
            <div
              v-for="(h, j) in [40, 70, 30, 85, 55, 90, 45, 65]"
              :key="j"
              :class="[kpi.barClass]"
              class="w-1 rounded-full"
              :style="{ height: h + '%' }"
            />
          </div>
        </div>
      </GlassCard>

      <!-- API Health -->
      <GlassCard :glow="store.health.ok ? 'glow-green' : 'glow-red'" class="kpi-card group relative overflow-hidden p-6">
        <div :class="['absolute inset-0 bg-gradient-to-br opacity-0 transition-opacity duration-500 group-hover:opacity-100', store.health.ok ? 'from-neon-green/20 to-transparent' : 'from-neon-red/20 to-transparent']" />
        <div v-if="store.loading" class="relative space-y-3">
          <SkeletonLoader :lines="2" />
        </div>
        <div v-else class="relative">
          <div class="flex items-center justify-between">
            <span class="text-xs font-medium uppercase tracking-wider text-gray-500">API Health</span>
            <span class="text-lg">♥</span>
          </div>
          <div class="mt-2 flex items-center gap-3">
            <span class="relative flex h-4 w-4">
              <span
                :class="[
                  'absolute inline-flex h-full w-full rounded-full opacity-50',
                  store.health.ok ? 'bg-neon-green animate-ping' : 'bg-neon-red animate-flicker',
                ]"
              />
              <span
                :class="[
                  'relative inline-flex h-4 w-4 rounded-full',
                  store.health.ok ? 'bg-neon-green' : 'bg-neon-red',
                ]"
              />
            </span>
            <span :class="['text-2xl font-bold', store.health.ok ? 'text-neon-green' : 'text-neon-red']">
              {{ store.health.ok ? "Healthy" : "Down" }}
            </span>
          </div>
          <!-- Uptime bar -->
          <div class="mt-3 h-1.5 w-full rounded-full bg-void">
            <div
              :class="[
                'h-full rounded-full transition-all duration-1000',
                store.health.ok ? 'bg-gradient-to-r from-neon-green to-neon-cyan w-full' : 'bg-neon-red w-1/4',
              ]"
            />
          </div>
        </div>
      </GlassCard>
    </div>

    <!-- System Info + Quick Actions -->
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <!-- Quick Actions -->
      <GlassCard class="p-6">
        <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">Quick Actions</h3>
        <div class="flex flex-wrap gap-3">
          <router-link to="/instances">
            <NeonButton variant="ghost">⬡ View Instances</NeonButton>
          </router-link>
          <NeonButton variant="ghost" @click="store.fetchInstances()" :loading="store.loading">
            ↻ Refresh Data
          </NeonButton>
        </div>
      </GlassCard>

      <!-- System Info -->
      <GlassCard class="p-6">
        <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">System</h3>
        <div class="space-y-2 font-mono text-xs text-gray-500">
          <div class="flex justify-between">
            <span>Runtime</span>
            <span class="text-neon-cyan">Kubernetes</span>
          </div>
          <div class="flex justify-between">
            <span>Database</span>
            <span class="text-neon-purple">CloudNativePG</span>
          </div>
          <div class="flex justify-between">
            <span>API</span>
            <span class="text-neon-blue">Go / Echo</span>
          </div>
          <div class="flex justify-between">
            <span>Auth</span>
            <span class="text-neon-green">JWT / RBAC</span>
          </div>
        </div>
      </GlassCard>
    </div>
  </div>
</template>
