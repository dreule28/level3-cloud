<!--
  AppShell — Cyberpunk layout with animated sidebar + system HUD topbar
-->
<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useInstanceStore } from "@/stores/instances";

const auth = useAuthStore();
const instanceStore = useInstanceStore();
const router = useRouter();
const route = useRoute();
const collapsed = ref(false);
const clock = ref("");
let clockInterval;

const navItems = [
  { path: "/dashboard", label: "Dashboard", icon: "◈" },
  { path: "/instances", label: "Instances", icon: "⬡" },
  { path: "/logs", label: "Logs", icon: "□" },
];

function updateClock() {
  const now = new Date();
  clock.value = now.toLocaleTimeString("en-US", { hour12: false });
}

onMounted(() => {
  updateClock();
  clockInterval = setInterval(updateClock, 1000);
});
onUnmounted(() => clearInterval(clockInterval));

function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <div class="flex h-screen overflow-hidden">
    <!-- ── Sidebar ── -->
    <aside
      :class="[
        'relative flex flex-col border-r border-border bg-surface/50 backdrop-blur-xl transition-all duration-300',
        collapsed ? 'w-16' : 'w-60',
      ]"
    >
      <!-- Sidebar glow accent -->
      <div class="absolute top-0 right-0 h-full w-px bg-gradient-to-b from-neon-blue/20 via-neon-purple/10 to-transparent" />

      <!-- Logo -->
      <div class="flex h-16 items-center gap-3 border-b border-border px-4">
        <div class="relative flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-neon-blue to-neon-purple text-base font-bold text-white shadow-[0_0_25px_rgba(139,92,246,0.3)]">
          ☁
          <div class="absolute inset-0 rounded-xl bg-gradient-to-br from-neon-blue to-neon-purple opacity-50 blur-md" />
        </div>
        <Transition name="fade">
          <div v-if="!collapsed" class="flex flex-col">
            <span class="text-sm font-bold tracking-wide text-gradient">CloudOS</span>
            <span class="text-[10px] font-mono text-gray-600">v1.0.0-cyberpunk</span>
          </div>
        </Transition>
      </div>

      <!-- Nav -->
      <nav class="mt-4 flex flex-1 flex-col gap-1 px-2">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          :class="[
            'group relative flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200',
            route.path.startsWith(item.path)
              ? 'bg-neon-blue/10 text-neon-blue'
              : 'text-gray-400 hover:bg-white/5 hover:text-white',
          ]"
        >
          <!-- Active indicator bar -->
          <div
            v-if="route.path.startsWith(item.path)"
            class="absolute left-0 top-1/2 h-6 w-0.5 -translate-y-1/2 rounded-full bg-neon-blue shadow-[0_0_8px_rgba(59,130,246,0.6)]"
          />
          <span class="text-lg transition-transform duration-200 group-hover:scale-110">{{ item.icon }}</span>
          <Transition name="fade">
            <span v-if="!collapsed">{{ item.label }}</span>
          </Transition>
        </router-link>
      </nav>

      <!-- Sidebar footer -->
      <div class="border-t border-border p-3">
        <button
          @click="collapsed = !collapsed"
          class="flex w-full items-center justify-center rounded-lg p-2 text-gray-500 transition hover:bg-white/5 hover:text-white"
        >
          <span :class="collapsed ? 'rotate-180' : ''" class="inline-block transition-transform duration-300">‹</span>
        </button>
      </div>
    </aside>

    <!-- ── Main Content ── -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <!-- Top bar / HUD -->
      <header class="flex h-16 items-center justify-between border-b border-border bg-surface/30 px-6 backdrop-blur-xl">
        <!-- Path display -->
        <div class="flex items-center gap-4">
          <div class="text-sm text-gray-500 font-mono">
            <span class="text-neon-cyan">~/</span>{{ route.path }}
          </div>
        </div>

        <!-- System indicators -->
        <div class="flex items-center gap-5">
          <!-- Status dot -->
          <div class="flex items-center gap-2">
            <span class="relative flex h-2 w-2">
              <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-neon-green opacity-50" />
              <span class="relative inline-flex h-2 w-2 rounded-full bg-neon-green" />
            </span>
            <span class="text-xs font-mono text-gray-500">ONLINE</span>
          </div>

          <!-- Instance count -->
          <div class="hidden sm:flex items-center gap-2 rounded-lg border border-border/50 bg-void/50 px-3 py-1">
            <span class="text-xs text-gray-500">⬡</span>
            <span class="text-xs font-mono text-neon-blue">{{ instanceStore.totalInstances }}</span>
          </div>

          <!-- Clock -->
          <span class="hidden sm:block text-xs font-mono text-gray-600">{{ clock }}</span>

          <!-- Logout -->
          <button
            @click="logout"
            class="rounded-lg border border-border px-4 py-1.5 text-xs font-medium text-gray-400 transition hover:border-neon-red/50 hover:text-neon-red hover:shadow-[0_0_15px_rgba(239,68,68,0.1)]"
          >
            Logout
          </button>
        </div>
      </header>

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-6">
        <slot />
      </main>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
