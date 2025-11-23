<template>
  <div class="card">
    <h2>Versions for {{ type }}</h2>

    <div class="row">
      <button @click="reload">Reload</button>
      <button @click="loadLatest">Show latest</button>
    </div>

    <div v-if="loading">Loading…</div>
    <div v-else-if="versions.length===0" class="small">No versions yet.</div>

    <ul>
      <li v-for="v in versions" :key="v.version" style="margin:8px 0;">
        <div><b>{{ v.version }}</b> ({{ formatSize(v.sizeBytes) }})</div>
        <div class="small">sha256: {{ v.sha256 }}</div>
        <div class="small">created: {{ formatDate(v.createdAt) }}</div>
        <div style="margin-top:4px;">
          <a :href="v.downloadUrl" target="_blank">Download</a>
          <button @click="remove(v.version)" style="margin-left:8px;">Delete</button>
        </div>
      </li>
    </ul>

    <div v-if="latest" class="small" style="margin-top:8px;">
      Latest: <b>{{ latest.version }}</b>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import { FirmwareAPI } from "../api";

const props = defineProps({ type: String });
const versions = ref([]);
const latest = ref(null);
const loading = ref(false);

watch(() => props.type, reload, { immediate: true });

async function reload() {
  if (!props.type) return;
  loading.value = true;
  latest.value = null;
  try {
    versions.value = await FirmwareAPI.list(props.type);
  } catch (e) {
    versions.value = [];
  } finally {
    loading.value = false;
  }
}

async function loadLatest() {
  latest.value = null;
  try {
    latest.value = await FirmwareAPI.latest(props.type);
  } catch {}
}

async function remove(version) {
  if (!confirm(`Delete ${props.type} ${version}?`)) return;
  await FirmwareAPI.remove(props.type, version);
  await reload();
}

function formatSize(n) {
  if (!n && n !== 0) return "-";
  const kb = n / 1024;
  if (kb < 1024) return kb.toFixed(1) + " KB";
  return (kb / 1024).toFixed(2) + " MB";
}

function formatDate(d) {
  try { return new Date(d).toLocaleString(); } catch { return d; }
}
</script>
