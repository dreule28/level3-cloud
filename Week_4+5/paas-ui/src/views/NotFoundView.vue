<!--
  NotFoundView — Animated glitch "Lost in the Cloud" 404 with GSAP + data corruption effect
-->
<script setup>
import { onMounted, onUnmounted, ref } from "vue";
import gsap from "gsap";
import NeonButton from "@/components/ui/NeonButton.vue";

const glitchText = ref("404");
const subText = ref("SYSTEM_ERROR::ROUTE_NOT_FOUND");
const chars = "!@#$%^&*()_+-=[]{}|;:',.<>?/~`0123456789ABCDEFabcdef";
const titleRef = ref(null);
let glitchInterval;

onMounted(() => {
  // GSAP entrance
  if (titleRef.value) {
    gsap.from(titleRef.value, {
      scale: 3,
      opacity: 0,
      rotation: -5,
      duration: 0.8,
      ease: "elastic.out(1, 0.5)",
    });
  }

  // Glitch effect
  glitchInterval = setInterval(() => {
    // Glitch the 404
    const glitched = "404"
      .split("")
      .map((ch) => (Math.random() > 0.6 ? chars[Math.floor(Math.random() * chars.length)] : ch))
      .join("");
    glitchText.value = glitched;
    setTimeout(() => (glitchText.value = "404"), 120);

    // Glitch subtitle too
    if (Math.random() > 0.7) {
      const orig = "SYSTEM_ERROR::ROUTE_NOT_FOUND";
      subText.value = orig
        .split("")
        .map((ch) => (Math.random() > 0.8 ? chars[Math.floor(Math.random() * chars.length)] : ch))
        .join("");
      setTimeout(() => (subText.value = orig), 150);
    }
  }, 1800);
});

onUnmounted(() => clearInterval(glitchInterval));
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center p-4 text-center">
    <!-- Glitch title -->
    <h1
      ref="titleRef"
      class="relative text-[8rem] font-black leading-none text-gradient select-none sm:text-[12rem]"
    >
      {{ glitchText }}
      <!-- Offset shadow layers for glitch -->
      <span class="absolute inset-0 text-neon-blue/20 translate-x-1 translate-y-0.5 blur-[1px] select-none">{{ glitchText }}</span>
      <span class="absolute inset-0 text-neon-red/15 -translate-x-1 -translate-y-0.5 blur-[1px] select-none">{{ glitchText }}</span>
    </h1>

    <!-- Subtitle -->
    <div class="mt-4 space-y-3">
      <p class="font-mono text-xs text-neon-red/70 tracking-widest animate-flicker">{{ subText }}</p>
      <p class="text-xl font-semibold text-white">Lost in the Cloud</p>
      <p class="text-sm text-gray-500">The resource you're looking for drifted into the void.</p>
    </div>

    <!-- Decorative line -->
    <div class="neon-line mx-auto mt-8 w-64" />

    <!-- CTA -->
    <div class="mt-8">
      <router-link to="/dashboard">
        <NeonButton>← Return to Dashboard</NeonButton>
      </router-link>
    </div>

    <!-- Ambient glow -->
    <div class="absolute top-1/4 h-[400px] w-[400px] rounded-full bg-neon-purple/10 blur-[150px] animate-glow-pulse" />
    <div class="absolute bottom-1/4 right-1/4 h-[300px] w-[300px] rounded-full bg-neon-red/8 blur-[120px] animate-float" />
  </div>
</template>
