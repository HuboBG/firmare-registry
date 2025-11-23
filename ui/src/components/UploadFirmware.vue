<template>
  <div class="card">
    <h2>Upload Firmware ({{ type }})</h2>

    <label>Version (semver)</label>
    <input v-model="version" placeholder="1.0.0"/>

    <input type="file" @change="onFile"/>
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

<script setup lang="ts">
import {ref} from "vue";
import {FirmwareAPI, type FirmwareDTO} from "../api";

const props = defineProps<{ type: string }>();
const emit = defineEmits<{ (e: "uploaded", dto: FirmwareDTO): void }>();

const version = ref<string>("");
const file = ref<File | null>(null);
const result = ref<FirmwareDTO | null>(null);
const error = ref<string>("");
const uploading = ref<boolean>(false);

function onFile(e: Event) {
  const input = e.target as HTMLInputElement;
  file.value = input.files?.[0] ?? null;
}

async function upload() {
  error.value = "";
  result.value = null;
  if (!file.value) return;

  uploading.value = true;
  try {
    const r = await FirmwareAPI.upload(props.type, version.value.trim(), file.value);
    result.value = r;
    emit("uploaded", r);
    version.value = "";
    file.value = null;
  } catch (e: any) {
    error.value = e?.response?.data || e?.message || "Upload failed";
  } finally {
    uploading.value = false;
  }
}
</script>
