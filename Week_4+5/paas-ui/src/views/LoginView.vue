<!--
  LoginView — Cyberpunk animated login with floating neon card + terminal typing effect
-->
<script setup>
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import NeonButton from "@/components/ui/NeonButton.vue";
import gsap from "gsap";

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();

const username = ref("");
const password = ref("");
const typedTitle = ref("");
const cardRef = ref(null);

const titleText = "CloudOS";

onMounted(() => {
  // Type-in effect for title
  let i = 0;
  const typeInterval = setInterval(() => {
    if (i < titleText.length) {
      typedTitle.value += titleText[i];
      i++;
    } else {
      clearInterval(typeInterval);
    }
  }, 120);

  // GSAP entrance for card
  if (cardRef.value) {
    gsap.from(cardRef.value, {
      y: 60,
      opacity: 0,
      scale: 0.92,
      duration: 0.8,
      ease: "power3.out",
      delay: 0.3,
    });
  }
});

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
    <div class="absolute top-1/3 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-neon-purple/10 blur-[120px] animate-glow-pulse" />
    <div class="absolute bottom-1/4 right-1/4 h-[300px] w-[300px] rounded-full bg-neon-blue/10 blur-[100px] animate-glow-pulse" style="animation-delay: 1s" />
    <div class="absolute top-1/4 left-1/4 h-[200px] w-[200px] rounded-full bg-neon-cyan/8 blur-[80px] animate-float" />

    <!-- Login Card -->
    <div ref="cardRef" class="relative w-full max-w-sm">
      <div class="rounded-2xl glass-strong p-8 glow-purple">
        <!-- Logo -->
        <div class="mb-8 text-center">
          <div class="mx-auto mb-4 relative flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-neon-blue to-neon-purple text-2xl shadow-[0_0_40px_rgba(139,92,246,0.4)]">
            ☁
            <div class="absolute inset-0 rounded-2xl bg-gradient-to-br from-neon-blue to-neon-purple opacity-40 blur-lg animate-glow-pulse" />
          </div>
          <h1 class="text-3xl font-bold text-gradient">
            {{ typedTitle }}<span class="animate-flicker text-neon-blue">_</span>
          </h1>
          <p class="mt-2 text-xs font-mono text-gray-600">PaaS Control Plane · v1.0.0</p>
        </div>

        <!-- Terminal-style auth prompt -->
        <div class="mb-5 rounded-lg border border-border/50 bg-void/60 px-3 py-2 font-mono text-[11px]">
          <span class="text-neon-green">root@cloudos</span><span class="text-gray-600">:</span><span class="text-neon-blue">~</span><span class="text-gray-500">$ </span><span class="text-gray-400">authenticate --user</span>
        </div>

        <!-- Form -->
        <form @submit.prevent="onLogin" class="space-y-5">
          <div>
            <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Username</label>
            <input
              v-model="username"
              type="text"
              autocomplete="username"
              class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 text-sm text-white placeholder-gray-600 transition-all duration-300 focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50 focus:shadow-[0_0_20px_rgba(59,130,246,0.15)]"
              placeholder="admin"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Password</label>
            <input
              v-model="password"
              type="password"
              autocomplete="current-password"
              class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 text-sm text-white placeholder-gray-600 transition-all duration-300 focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50 focus:shadow-[0_0_20px_rgba(59,130,246,0.15)]"
              placeholder="••••••••"
            />
          </div>

          <!-- Error -->
          <Transition name="slide">
            <p v-if="auth.error" class="rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-2 text-xs text-neon-red animate-flicker">
              {{ auth.error }}
            </p>
          </Transition>

          <NeonButton type="submit" :loading="auth.loading" class="w-full justify-center">
            Sign In
          </NeonButton>
        </form>

        <!-- Bottom glow line -->
        <div class="neon-line mt-6" />
        <p class="mt-4 text-center text-[10px] font-mono text-gray-700">
          JWT · TLS · E2E Encrypted · Kubernetes
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slide-enter-active, .slide-leave-active { transition: all 0.3s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
