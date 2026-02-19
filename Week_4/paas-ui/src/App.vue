<!--
  App.vue — Root component
  Routes through layout (blank for login/404, shell for authenticated pages)
  GSAP-powered page transitions + ambient effects layer
-->
<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import gsap from "gsap";
import AppShell from "@/components/layout/AppShell.vue";
import BackgroundGrid from "@/components/effects/BackgroundGrid.vue";
import CursorGlow from "@/components/effects/CursorGlow.vue";
import ToastContainer from "@/components/ui/ToastContainer.vue";

const route = useRoute();
const isBlankLayout = computed(() => route.meta.layout === "blank");

/* GSAP page transitions */
function onBeforeEnter(el) {
  gsap.set(el, { opacity: 0, y: 20, scale: 0.98 });
}
function onEnter(el, done) {
  gsap.to(el, { opacity: 1, y: 0, scale: 1, duration: 0.45, ease: "power3.out", onComplete: done });
}
function onLeave(el, done) {
  gsap.to(el, { opacity: 0, y: -12, scale: 0.99, duration: 0.25, ease: "power2.in", onComplete: done });
}
</script>

<template>
  <BackgroundGrid />
  <CursorGlow />
  <ToastContainer />
  <div class="relative z-10 min-h-screen">
    <!-- Blank layout: login, 404 -->
    <template v-if="isBlankLayout">
      <router-view v-slot="{ Component }">
        <Transition
          :css="false"
          @before-enter="onBeforeEnter"
          @enter="onEnter"
          @leave="onLeave"
          mode="out-in"
        >
          <component :is="Component" />
        </Transition>
      </router-view>
    </template>
    <!-- Shell layout: sidebar + main content -->
    <AppShell v-else>
      <router-view v-slot="{ Component }">
        <Transition
          :css="false"
          @before-enter="onBeforeEnter"
          @enter="onEnter"
          @leave="onLeave"
          mode="out-in"
        >
          <component :is="Component" />
        </Transition>
      </router-view>
    </AppShell>
  </div>
</template>
