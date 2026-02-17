import { createRouter, createWebHistory } from "vue-router";
import LoginView from "./views/LoginView.vue";
import InstancesView from "./views/InstancesView.vue";
import { getToken } from "./api/http";

const routes = [
  { path: "/login", component: LoginView },
  { path: "/", redirect: "/instances" },
  { path: "/instances", component: InstancesView, meta: { requiresAuth: true } },
  { path: "/instances/:id", component: () => import("./views/InstanceDetailView.vue"), meta: { requiresAuth: true } }

];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getToken()) {
    return "/login";
  }
  if (to.path === "/login" && getToken()) {
    return "/instances";
  }
});

export default router;
