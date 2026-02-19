/**
 * CloudOS — PaaS Control Plane
 * App entry: Vue 3 + Pinia + Router + Global CSS
 */
import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";
import "./style.css";

const app = createApp(App);
app.use(createPinia());
app.use(router);
app.mount("#app");
