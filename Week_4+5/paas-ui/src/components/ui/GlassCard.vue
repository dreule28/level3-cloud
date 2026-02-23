<!--
  GlassCard — Glassmorphic container with animated gradient border on hover
-->
<script setup>
defineProps({
  glow: { type: String, default: "" },
  hoverable: { type: Boolean, default: true },
});
</script>

<template>
  <div
    :class="[
      'group relative rounded-xl glass transition-all duration-300',
      hoverable && 'gradient-border hover:bg-glass-hover hover:border-white/10 hover:-translate-y-0.5',
      glow,
    ]"
  >
    <slot />
  </div>
</template>
