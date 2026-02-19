<!--
  NotFoundView — Animated glitch "Lost in the Cloud" 404
-->
<script setup>
import { onMounted, ref } from "vue";
import NeonButton from "@/components/ui/NeonButton.vue";

const glitchText = ref("404");
const chars = "!@#$%^&*()_+-=[]{}|;:',.<>?/~`";
let interval;

onMounted(() => {
  interval = setInterval(() => {
    const glitched = "404"
      .split("")
      .map((ch) => (Math.random() > 0.7 ? chars[Math.floor(Math.random() * chars.length)] : ch))
      .join("");
    glitchText.value = glitched;
    setTimeout(() => (glitchText.value = "404"), 100);
  }, 2000);
});
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center p-4 text-center">
    <!-- Glitch title -->
    <h1 class="text-[8rem] font-black leading-none text-gradient animate-flicker select-none sm:text-[12rem]">
      {{ glitchText }}
    </h1>

    <!-- Subtitle -->
    <div class="mt-4 space-y-2">
      <p class="text-xl font-semibold text-white">Lost in the Cloud</p>
      <p class="text-sm text-gray-500">The resource you're looking for drifted into the void.</p>
    </div>

    <!-- Decorative line -->
    <div class="neon-line mx-auto mt-8 w-48" />

    <!-- CTA -->
    <div class="mt-8">
      <router-link to="/dashboard">
        <NeonButton>← Return to Dashboard</NeonButton>
      </router-link>
    </div>

    <!-- Ambient glow -->
    <div class="absolute top-1/4 h-[400px] w-[400px] rounded-full bg-neon-purple/10 blur-[150px]" />
    <div class="absolute bottom-1/4 right-1/4 h-[300px] w-[300px] rounded-full bg-neon-red/8 blur-[120px]" />
  </div>
</template>
