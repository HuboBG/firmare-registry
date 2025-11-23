<template>
  <div class="app">
    <header class="header">
      <h1>Firmware Registry Admin</h1>
      <nav>
        <button :class="{active:tab==='firmware'}" @click="tab='firmware'">Firmware</button>
        <button :class="{active:tab==='webhooks'}" @click="tab='webhooks'">Webhooks</button>
        <button :class="{active:tab==='settings'}" @click="tab='settings'">Settings</button>
      </nav>
    </header>

    <main>
      <FirmwareTypes
          v-if="tab==='firmware'"
          @selectType="selectType"
      />

      <div v-if="tab==='firmware' && selectedType" class="twoCol">
        <FirmwareVersions :type="selectedType"/>
        <UploadFirmware :type="selectedType" @uploaded="onUploaded"/>
      </div>

      <Webhooks v-if="tab==='webhooks'"/>
      <Settings v-if="tab==='settings'"/>
    </main>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import FirmwareTypes from "./components/FirmwareTypes.vue";
import FirmwareVersions from "./components/FirmwareVersions.vue";
import UploadFirmware from "./components/UploadFirmware.vue";
import Webhooks from "./components/Webhooks.vue";
import Settings from "./components/Settings.vue";
import type {FirmwareDTO} from "./api";

const tab = ref<"firmware" | "webhooks" | "settings">("firmware");
const selectedType = ref<string>("");

// emitted from FirmwareTypes
function selectType(t: string) {
  selectedType.value = t;
}

// currently unused, but handy for future refresh wiring
function onUploaded(_dto: FirmwareDTO) {
  // no-op for now
}
</script>

<style>
.app {
  font-family: system-ui, sans-serif;
  padding: 16px;
  max-width: 1100px;
  margin: 0 auto;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

nav button {
  margin-left: 8px;
}

nav button.active {
  font-weight: 700;
}

.twoCol {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 12px;
}

.card {
  border: 1px solid #ddd;
  border-radius: 12px;
  padding: 12px;
  background: #fff;
}

input, button, select {
  padding: 6px 8px;
  margin: 4px 0;
}

.small {
  font-size: 12px;
  color: #555;
}

.row {
  display: flex;
  gap: 8px;
  align-items: center;
}

ul {
  padding-left: 16px;
}
</style>
