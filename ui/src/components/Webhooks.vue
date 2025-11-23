<template>
  <div class="card">
    <h2>Webhooks</h2>
    <button @click="reload">Reload</button>

    <h3 style="margin-top:12px;">Add webhook</h3>
    <input v-model="url" placeholder="https://yourapp/webhook/fw" />

    <div class="row">
      <label><input type="checkbox" v-model="evtUploaded" /> firmware.uploaded</label>
      <label><input type="checkbox" v-model="evtDeleted" /> firmware.deleted</label>
      <label><input type="checkbox" v-model="enabled" /> enabled</label>
    </div>

    <button @click="create">Create</button>

    <hr />

    <ul v-if="hooks.length">
      <li v-for="h in hooks" :key="h.id" style="margin:8px 0;">
        <div><b>#{{ h.id }}</b> {{ h.url }}</div>
        <div class="small">events: {{ h.events.join(", ") }}</div>
        <div class="small">enabled: {{ h.enabled }}</div>
        <button @click="remove(h.id)">Delete</button>
      </li>
    </ul>

    <p v-else class="small">No webhooks registered.</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { WebhookAPI } from "../api";

const hooks = ref([]);

const url = ref("");
const evtUploaded = ref(true);
const evtDeleted = ref(false);
const enabled = ref(true);

onMounted(reload);

async function reload() {
  hooks.value = await WebhookAPI.list();
}

async function create() {
  const events = [];
  if (evtUploaded.value) events.push("firmware.uploaded");
  if (evtDeleted.value) events.push("firmware.deleted");
  if (!url.value.trim() || events.length === 0) return;

  await WebhookAPI.create({
    url: url.value.trim(),
    events,
    enabled: enabled.value
  });
  url.value = "";
  await reload();
}

async function remove(id) {
  if (!confirm(`Delete webhook #${id}?`)) return;
  await WebhookAPI.remove(id);
  await reload();
}
</script>
