<!--
  NeonButton — Animated gradient button with ripple + glow
-->
<script setup>
import { ref } from "vue";

defineProps({
  variant: { type: String, default: "primary" },
  loading: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
});

const ripple = ref(false);
function doRipple() {
  ripple.value = true;
  setTimeout(() => (ripple.value = false), 600);
}
</script>

<template>
  <button
    @click="doRipple"
    :disabled="disabled || loading"
    :class="[
      'relative overflow-hidden rounded-lg px-5 py-2.5 text-sm font-semibold transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-void disabled:opacity-40 disabled:cursor-not-allowed',
      variant === 'primary' && 'bg-gradient-to-r from-neon-blue to-neon-purple text-white hover:shadow-[0_0_30px_rgba(59,130,246,0.3)] focus:ring-neon-blue',
      variant === 'danger' && 'bg-gradient-to-r from-neon-red to-neon-pink text-white hover:shadow-[0_0_30px_rgba(239,68,68,0.3)] focus:ring-neon-red',
      variant === 'ghost' && 'border border-border text-gray-400 hover:border-neon-blue/50 hover:text-white hover:bg-white/5 focus:ring-neon-blue',
      variant === 'success' && 'bg-gradient-to-r from-neon-green to-neon-cyan text-white hover:shadow-[0_0_30px_rgba(16,185,129,0.3)] focus:ring-neon-green',
    ]"
  >
    <!-- Ripple -->
    <span
      v-if="ripple"
      class="absolute inset-0 animate-ping rounded-lg bg-white/10"
    />

    <!-- Loading spinner -->
    <span v-if="loading" class="mr-2 inline-block h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" />

    <slot />
  </button>
</template>
