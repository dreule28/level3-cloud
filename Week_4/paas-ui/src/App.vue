<!--
  App.vue — Root component
  Routes through layout (blank for login/404, shell for authenticated pages)
-->
<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import AppShell from "@/components/layout/AppShell.vue";
import BackgroundGrid from "@/components/effects/BackgroundGrid.vue";
import CursorGlow from "@/components/effects/CursorGlow.vue";

const route = useRoute();
const isBlankLayout = computed(() => route.meta.layout === "blank");
</script>

<template>
  <BackgroundGrid />
  <CursorGlow />
  <div class="relative z-10 min-h-screen">
    <!-- Blank layout: login, 404 -->
    <router-view v-if="isBlankLayout" />
    <!-- Shell layout: sidebar + main content -->
    <AppShell v-else>
      <router-view v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </router-view>
    </AppShell>
  </div>
</template>

<style>
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}
.page-enter-from {
  opacity: 0;
  transform: translateY(12px);
}
.page-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
