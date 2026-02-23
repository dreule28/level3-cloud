/**
 * Vue Router — CloudOS
 * Auth guards, lazy-loaded views, page transitions
 */
import { createRouter, createWebHistory } from "vue-router";
import { getToken } from "@/api/auth";

const routes = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginView.vue"),
    meta: { layout: "blank" },
  },
  {
    path: "/",
    redirect: "/dashboard",
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: () => import("@/views/DashboardView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/instances",
    name: "instances",
    component: () => import("@/views/InstancesView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/instances/:id",
    name: "instance-detail",
    component: () => import("@/views/InstanceDetailView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "not-found",
    component: () => import("@/views/NotFoundView.vue"),
    meta: { layout: "blank" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Auth guard
router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getToken()) {
    return { path: "/login", query: { redirect: to.fullPath } };
  }
  if (to.path === "/login" && getToken()) {
    return "/dashboard";
  }
});

export default router;
