<!--
  CursorGlow — Soft neon glow follows the mouse
-->
<script setup>
import { ref, onMounted, onUnmounted } from "vue";

const x = ref(-100);
const y = ref(-100);

function onMove(e) {
  x.value = e.clientX;
  y.value = e.clientY;
}

onMounted(() => window.addEventListener("mousemove", onMove));
onUnmounted(() => window.removeEventListener("mousemove", onMove));
</script>

<template>
  <div
    class="pointer-events-none fixed z-50 h-[300px] w-[300px] rounded-full opacity-20 blur-[100px] transition-transform duration-100 ease-out"
    :style="{
      background: 'radial-gradient(circle, rgba(59,130,246,0.4), rgba(139,92,246,0.2), transparent)',
      transform: `translate(${x - 150}px, ${y - 150}px)`,
    }"
  />
</template>
