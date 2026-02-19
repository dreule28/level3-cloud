<!--
  StatusPill — Animated status indicator
  READY = neon green glow, CREATING = pulse, ERROR = red flicker
-->
<script setup>
import { computed } from "vue";
const props = defineProps({ status: { type: String, required: true } });

const normalized = computed(() => props.status?.toLowerCase() || "unknown");

const config = computed(() => {
  switch (normalized.value) {
    case "ready":
      return { label: "Ready", class: "bg-neon-green/15 text-neon-green border-neon-green/30", dot: "bg-neon-green", glow: "glow-green", anim: "" };
    case "creating":
      return { label: "Creating", class: "bg-neon-amber/15 text-neon-amber border-neon-amber/30", dot: "bg-neon-amber", glow: "", anim: "animate-glow-pulse" };
    case "error":
      return { label: "Error", class: "bg-neon-red/15 text-neon-red border-neon-red/30", dot: "bg-neon-red", glow: "glow-red", anim: "animate-flicker" };
    default:
      return { label: props.status, class: "bg-gray-500/15 text-gray-400 border-gray-500/30", dot: "bg-gray-400", glow: "", anim: "" };
  }
});
</script>

<template>
  <span
    :class="[
      'inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-medium',
      config.class,
      config.glow,
      config.anim,
    ]"
  >
    <span :class="['h-1.5 w-1.5 rounded-full', config.dot, config.anim ? 'animate-glow-pulse' : '']" />
    {{ config.label }}
  </span>
</template>
