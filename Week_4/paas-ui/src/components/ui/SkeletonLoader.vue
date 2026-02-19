<!--
  SkeletonLoader — Animated shimmer skeleton for loading states
-->
<script setup>
defineProps({
  lines: { type: Number, default: 3 },
  height: { type: String, default: "h-4" },
});
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="i in lines"
      :key="i"
      :class="[
        'rounded-lg bg-gradient-to-r from-surface via-surface-light to-surface bg-[length:200%_100%] animate-shimmer',
        height,
        i === lines ? 'w-2/3' : 'w-full',
      ]"
    />
  </div>
</template>
