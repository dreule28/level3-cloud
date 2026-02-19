<!--
  ToastContainer — Renders stacked toast notifications with slide animations
-->
<script setup>
import { useToast } from "@/composables/useToast";

const { toasts, dismiss } = useToast();

const iconMap = {
  success: "✓",
  error: "✕",
  info: "ℹ",
};

const colorMap = {
  success: "border-neon-green/40 bg-neon-green/10 text-neon-green",
  error: "border-neon-red/40 bg-neon-red/10 text-neon-red",
  info: "border-neon-blue/40 bg-neon-blue/10 text-neon-blue",
};

const glowMap = {
  success: "shadow-[0_0_20px_rgba(16,185,129,0.15)]",
  error: "shadow-[0_0_20px_rgba(239,68,68,0.15)]",
  info: "shadow-[0_0_20px_rgba(59,130,246,0.15)]",
};
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="[
            'pointer-events-auto flex items-center gap-3 rounded-xl border px-4 py-3 text-sm font-medium backdrop-blur-xl cursor-pointer transition-all duration-300',
            colorMap[toast.type] || colorMap.info,
            glowMap[toast.type] || glowMap.info,
            toast.leaving ? 'opacity-0 translate-x-8' : '',
          ]"
          @click="dismiss(toast.id)"
        >
          <span class="text-base">{{ iconMap[toast.type] || "ℹ" }}</span>
          <span>{{ toast.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active { transition: all 0.35s cubic-bezier(0.21, 1.02, 0.73, 1); }
.toast-leave-active { transition: all 0.25s ease-in; }
.toast-enter-from { opacity: 0; transform: translateX(80px) scale(0.9); }
.toast-leave-to { opacity: 0; transform: translateX(80px) scale(0.9); }
</style>
