/**
 * useToast — Global toast notification composable
 * Cyberpunk-styled toast with auto-dismiss + stacking
 */
import { ref } from "vue";

const toasts = ref([]);
let nextId = 0;

export function useToast() {
  function add(message, type = "info", duration = 4000) {
    const id = nextId++;
    toasts.value.push({ id, message, type, leaving: false });
    if (duration > 0) {
      setTimeout(() => dismiss(id), duration);
    }
  }

  function dismiss(id) {
    const t = toasts.value.find((t) => t.id === id);
    if (t) t.leaving = true;
    setTimeout(() => {
      toasts.value = toasts.value.filter((t) => t.id !== id);
    }, 300);
  }

  return {
    toasts,
    success: (msg) => add(msg, "success", 3000),
    error: (msg) => add(msg, "error", 5000),
    info: (msg) => add(msg, "info", 4000),
    dismiss,
  };
}
