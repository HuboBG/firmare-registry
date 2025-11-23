<template>
  <div class="card">
    <h2>Upload Firmware ({{ type }})</h2>

    <label>Version (semver)</label>
    <input v-model="version" placeholder="1.0.0" />

    <input type="file" @change="onFile" />
    <button :disabled="!file || !version || uploading" @click="upload">
      {{ uploading ? "Uploading…" : "Upload" }}
    </button>

    <div v-if="result" class="small" style="margin-top:8px;">
      Uploaded {{ result.type }} {{ result.version }} (sha={{ result.sha256 }})
    </div>

    <div v-if="error" class="small" style="margin-top:8px; color:#b00;">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { FirmwareAPI } from "../api";

const props = defineProps({ type: String });
const emit = defineEmits(["uploaded"]);

const version = ref("");
const file = ref(null);
const result = ref(null);
const error = ref("");
const uploading = ref(false);

function onFile(e) {
  file.value = e.target.files?.[0] || null;
}

async function upload() {
  error.value = "";
  result.value = null;
  uploading.value = true;
  try {
    const r = await FirmwareAPI.upload(props.type, version.value.trim(), file.value);
    result.value = r;
    emit("uploaded", r);
    version.value = "";
    file.value = null;
  } catch (e) {
    error.value = e?.response?.data || e?.message || "Upload failed";
  } finally {
    uploading.value = false;
  }
}
</script>
