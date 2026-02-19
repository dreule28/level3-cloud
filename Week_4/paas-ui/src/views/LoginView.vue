<!--
  LoginView — Cyberpunk animated login with floating neon card
-->
<script setup>
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import NeonButton from "@/components/ui/NeonButton.vue";

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();

const username = ref("");
const password = ref("");

async function onLogin() {
  if (!username.value || !password.value) return;
  const ok = await auth.login(username.value, password.value);
  if (ok) {
    const redirect = route.query.redirect || "/dashboard";
    router.push(redirect);
  }
}
</script>

<template>
  <div class="relative flex min-h-screen items-center justify-center p-4">
    <!-- Extra ambient glow behind form -->
    <div class="absolute top-1/3 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-neon-purple/10 blur-[120px]" />
    <div class="absolute bottom-1/4 right-1/4 h-[300px] w-[300px] rounded-full bg-neon-blue/10 blur-[100px]" />

    <!-- Login Card -->
    <div class="relative w-full max-w-sm animate-float">
      <div class="rounded-2xl glass-strong p-8 glow-purple">
        <!-- Logo -->
        <div class="mb-8 text-center">
          <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-neon-blue to-neon-purple text-2xl shadow-[0_0_40px_rgba(139,92,246,0.3)]">
            ☁
          </div>
          <h1 class="text-2xl font-bold text-gradient">CloudOS</h1>
          <p class="mt-1 text-sm text-gray-500">PaaS Control Plane</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="onLogin" class="space-y-5">
          <div>
            <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Username</label>
            <input
              v-model="username"
              type="text"
              autocomplete="username"
              class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 text-sm text-white placeholder-gray-600 transition focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50 focus:shadow-[0_0_20px_rgba(59,130,246,0.15)]"
              placeholder="admin"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Password</label>
            <input
              v-model="password"
              type="password"
              autocomplete="current-password"
              class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 text-sm text-white placeholder-gray-600 transition focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50 focus:shadow-[0_0_20px_rgba(59,130,246,0.15)]"
              placeholder="••••••••"
            />
          </div>

          <!-- Error -->
          <Transition name="slide">
            <p v-if="auth.error" class="rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-2 text-xs text-neon-red">
              {{ auth.error }}
            </p>
          </Transition>

          <NeonButton type="submit" :loading="auth.loading" class="w-full justify-center">
            Sign In
          </NeonButton>
        </form>

        <!-- Bottom glow line -->
        <div class="neon-line mt-6" />
        <p class="mt-4 text-center text-xs text-gray-600">
          Secured with JWT · E2E encrypted
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-enter-active, .slide-leave-active { transition: all 0.3s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
