<!--
  AppShell — Main layout with animated sidebar + top bar
-->
<script setup>
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const collapsed = ref(false);

const navItems = [
  { path: "/dashboard", label: "Dashboard", icon: "◈" },
  { path: "/instances", label: "Instances", icon: "⬡" },
];

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
      <!-- Logo -->
      <div class="flex h-16 items-center gap-3 border-b border-border px-4">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-neon-blue to-neon-purple text-xs font-bold text-white">
          ☁
        </div>
        <Transition name="fade">
          <span v-if="!collapsed" class="text-sm font-semibold tracking-wide text-gradient">CloudOS</span>
        </Transition>
      </div>

      <!-- Nav -->
      <nav class="mt-4 flex flex-1 flex-col gap-1 px-2">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          :class="[
            'group flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200',
            route.path.startsWith(item.path)
              ? 'bg-neon-blue/10 text-neon-blue glow-blue'
              : 'text-gray-400 hover:bg-white/5 hover:text-white',
          ]"
        >
          <span class="text-lg">{{ item.icon }}</span>
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
          <span :class="collapsed ? 'rotate-180' : ''" class="transition-transform duration-300">‹</span>
        </button>
      </div>
    </aside>

    <!-- ── Main Content ── -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <!-- Top bar -->
      <header class="flex h-16 items-center justify-between border-b border-border bg-surface/30 px-6 backdrop-blur-xl">
        <div class="text-sm text-gray-500 font-mono">
          <span class="text-neon-cyan">~/</span>{{ route.path }}
        </div>
        <button
          @click="logout"
          class="rounded-lg border border-border px-4 py-1.5 text-xs font-medium text-gray-400 transition hover:border-neon-red/50 hover:text-neon-red"
        >
          Logout
        </button>
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
