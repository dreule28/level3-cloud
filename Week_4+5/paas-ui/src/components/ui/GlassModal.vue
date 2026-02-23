<!--
  GlassModal — Glassmorphic modal with backdrop + animated entry
-->
<script setup>
defineProps({
  show: { type: Boolean, required: true },
  title: { type: String, default: "" },
});
defineEmits(["close"]);
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')" />
        <!-- Content -->
        <div class="relative w-full max-w-md animate-scale-in rounded-2xl glass-strong p-6 glow-purple">
          <div v-if="title" class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
            <button @click="$emit('close')" class="text-gray-500 transition hover:text-white">✕</button>
          </div>
          <slot />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.25s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
