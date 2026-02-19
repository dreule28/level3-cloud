<!--
  InstancesView — Cyberpunk animated table + GSAP stagger + create/delete modals + toast
-->
<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useInstanceStore } from "@/stores/instances";
import { useToast } from "@/composables/useToast";
import GlassCard from "@/components/ui/GlassCard.vue";
import NeonButton from "@/components/ui/NeonButton.vue";
import StatusPill from "@/components/ui/StatusPill.vue";
import GlassModal from "@/components/ui/GlassModal.vue";
import SkeletonLoader from "@/components/ui/SkeletonLoader.vue";

const store = useInstanceStore();
const router = useRouter();
const toast = useToast();

// Create modal
const showCreate = ref(false);
const newId = ref("");
const newInstances = ref(1);
const newStorage = ref(10);
const creating = ref(false);

// Delete modal
const showDelete = ref(false);
const deleteTarget = ref(null);
const deleting = ref(false);

onMounted(() => store.fetchInstances());

async function onCreate() {
  if (!newId.value.trim()) return;
  creating.value = true;
  const ok = await store.addInstance(newId.value.trim(), newInstances.value, newStorage.value);
  creating.value = false;
  if (ok) {
    toast.success(`Instance "${newId.value}" created`);
    newId.value = "";
    newInstances.value = 1;
    newStorage.value = 10;
    showCreate.value = false;
  } else {
    toast.error(store.error || "Failed to create instance");
  }
}

function confirmDelete(id) {
  deleteTarget.value = id;
  showDelete.value = true;
}

async function onDelete() {
  if (!deleteTarget.value) return;
  deleting.value = true;
  const ok = await store.removeInstance(deleteTarget.value);
  deleting.value = false;
  if (ok) {
    toast.success(`Instance "${deleteTarget.value}" deleted`);
  } else {
    toast.error(store.error || "Failed to delete instance");
  }
  showDelete.value = false;
  deleteTarget.value = null;
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-white">Instances</h1>
        <p class="mt-1 text-sm text-gray-500 font-mono">
          <span class="text-neon-cyan">db</span>.clusters — Manage database instances
        </p>
      </div>
      <div class="flex gap-3">
        <NeonButton variant="ghost" @click="store.fetchInstances" :loading="store.loading">
          ↻ Refresh
        </NeonButton>
        <NeonButton @click="showCreate = true">
          + Create Instance
        </NeonButton>
      </div>
    </div>

    <!-- Error banner -->
    <Transition name="slide">
      <div v-if="store.error" class="rounded-lg border border-neon-red/30 bg-neon-red/10 px-4 py-3 text-sm text-neon-red flex items-center justify-between">
        <span>{{ store.error }}</span>
        <button @click="store.clearError" class="text-neon-red/70 hover:text-neon-red underline text-xs">dismiss</button>
      </div>
    </Transition>

    <!-- Table -->
    <GlassCard :hoverable="false" class="overflow-hidden">
      <!-- Header -->
      <div class="grid grid-cols-12 gap-4 border-b border-border px-6 py-3 text-xs font-medium uppercase tracking-wider text-gray-600">
        <div class="col-span-1">#</div>
        <div class="col-span-4">Instance ID</div>
        <div class="col-span-3">Status</div>
        <div class="col-span-4 text-right">Actions</div>
      </div>

      <!-- Loading skeleton -->
      <div v-if="store.loading && store.instances.length === 0" class="p-6 space-y-4">
        <SkeletonLoader v-for="i in 3" :key="i" :lines="1" height="h-12" />
      </div>

      <!-- Rows -->
      <div>
        <div
          v-for="(inst, i) in store.instances"
          :key="inst.id"
          class="animate-fade-in group grid grid-cols-12 items-center gap-4 border-b border-border/30 px-6 py-4 transition-all duration-200 hover:bg-white/[0.03] hover:border-neon-blue/10"
          :style="{ animationDelay: `${i * 60}ms` }"
        >
          <div class="col-span-1 text-xs font-mono text-gray-700">
            {{ String(i + 1).padStart(2, "0") }}
          </div>
          <div class="col-span-4">
            <button
              @click="router.push(`/instances/${inst.id}`)"
              class="font-mono text-sm text-neon-cyan transition-all duration-200 hover:text-white hover:tracking-wider group-hover:text-shadow-[0_0_10px_rgba(6,182,212,0.5)]"
            >
              {{ inst.id }}
            </button>
          </div>
          <div class="col-span-3">
            <StatusPill :status="inst.status" />
          </div>
          <div class="col-span-4 flex justify-end gap-2">
            <NeonButton variant="ghost" @click="router.push(`/instances/${inst.id}`)">
              Details
            </NeonButton>
            <NeonButton variant="danger" @click="confirmDelete(inst.id)">
              Delete
            </NeonButton>
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="!store.loading && store.instances.length === 0" class="flex flex-col items-center justify-center py-20 text-gray-500">
        <span class="mb-4 text-5xl animate-float">⬡</span>
        <p class="text-sm font-medium">No instances deployed</p>
        <p class="mt-1 text-xs text-gray-600">Spin up your first database cluster</p>
        <NeonButton class="mt-6" @click="showCreate = true">+ Create Instance</NeonButton>
      </div>
    </GlassCard>

    <!-- ── Create Modal ── -->
    <GlassModal :show="showCreate" title="Create Instance" @close="showCreate = false">
      <form @submit.prevent="onCreate" class="space-y-4">
        <div>
          <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Instance ID</label>
          <input
            v-model="newId"
            type="text"
            placeholder="e.g. pg-production"
            class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 font-mono text-sm text-white placeholder-gray-600 transition-all duration-300 focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50 focus:shadow-[0_0_15px_rgba(59,130,246,0.1)]"
          />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Replicas</label>
            <input
              v-model.number="newInstances"
              type="number"
              min="1"
              max="5"
              class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 text-sm text-white transition-all duration-300 focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium uppercase tracking-wider text-gray-500">Storage (Gi)</label>
            <input
              v-model.number="newStorage"
              type="number"
              min="1"
              class="w-full rounded-lg border border-border bg-void/50 px-4 py-3 text-sm text-white transition-all duration-300 focus:border-neon-blue focus:outline-none focus:ring-1 focus:ring-neon-blue/50"
            />
          </div>
        </div>

        <Transition name="slide">
          <p v-if="store.error" class="rounded-lg border border-neon-red/30 bg-neon-red/10 px-3 py-2 text-xs text-neon-red">
            {{ store.error }}
          </p>
        </Transition>

        <div class="flex justify-end gap-3 pt-2">
          <NeonButton variant="ghost" type="button" @click="showCreate = false">Cancel</NeonButton>
          <NeonButton type="submit" :loading="creating">Deploy</NeonButton>
        </div>
      </form>
    </GlassModal>

    <!-- ── Delete Confirmation ── -->
    <GlassModal :show="showDelete" title="⚠ Confirm Deletion" @close="showDelete = false">
      <p class="text-sm text-gray-400">
        This will permanently destroy instance
        <span class="font-mono text-neon-red font-bold">{{ deleteTarget }}</span>
        and all associated data. This action cannot be undone.
      </p>
      <div class="mt-3 rounded-lg border border-neon-red/20 bg-neon-red/5 px-3 py-2 font-mono text-[11px] text-neon-red/70">
        kubectl delete cluster {{ deleteTarget }} --cascade=foreground
      </div>
      <div class="mt-6 flex justify-end gap-3">
        <NeonButton variant="ghost" @click="showDelete = false">Cancel</NeonButton>
        <NeonButton variant="danger" :loading="deleting" @click="onDelete">Destroy</NeonButton>
      </div>
    </GlassModal>
  </div>
</template>

<style scoped>
.list-enter-active { transition: all 0.4s ease; }
.list-leave-active { transition: all 0.3s ease; }
.list-enter-from { opacity: 0; transform: translateX(-20px); }
.list-leave-to { opacity: 0; transform: translateX(20px); }
.slide-enter-active, .slide-leave-active { transition: all 0.3s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
