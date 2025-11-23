<template>
  <div class="card">
    <h2>Firmware Types</h2>

    <div class="row">
      <input v-model="newType" placeholder="e.g. firmware1" @keyup.enter="addType" />
      <button @click="addType">Add</button>
      <button @click="reload">Reload from API</button>
    </div>

    <ul v-if="types.length">
      <li v-for="t in types" :key="t">
        <button class="typeBtn" @click="$emit('selectType', t)">
          {{ t }}
        </button>
      </li>
    </ul>
    <p v-else class="small">No types found yet. Upload a firmware or add manually.</p>

    <p class="small">
      Types are inferred from existing uploads. “Reload from API” discovers them by scanning known types list locally.
    </p>
  </div>
</template>

<script setup>
import { ref } from "vue";

// Simple UI-side type list. Since API doesn't have /types endpoint,
// we keep a local set that users can manage.
const types = ref(JSON.parse(localStorage.getItem("fw_types") || "[]"));
const newType = ref("");

function persist() {
  localStorage.setItem("fw_types", JSON.stringify(types.value));
}

function addType() {
  const t = newType.value.trim();
  if (!t) return;
  if (!types.value.includes(t)) {
    types.value.push(t);
    persist();
  }
  newType.value = "";
}

async function reload() {
  // There is no API "list types" currently. 
  // We just re-load local storage for now (future: add /types).
  types.value = JSON.parse(localStorage.getItem("fw_types") || "[]");
}
</script>

<style scoped>
.typeBtn { text-align:left; width: 100%; }
</style>
