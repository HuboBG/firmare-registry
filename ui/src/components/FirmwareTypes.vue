<template>
  <div class="card">
    <h2>Firmware Types</h2>

    <div class="row">
      <input v-model="newType" placeholder="e.g. firmware1" @keyup.enter="addType"/>
      <button @click="addType">Add</button>
      <button @click="reload">Reload</button>
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
      Types are stored locally (localStorage). API doesn't yet expose a types endpoint.
    </p>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";

// const emit = defineEmits<{ (e: "selectType", t: string): void }>();

const types = ref<string[]>(JSON.parse(localStorage.getItem("fw_types") || "[]"));
const newType = ref<string>("");

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

function reload() {
  types.value = JSON.parse(localStorage.getItem("fw_types") || "[]");
}
</script>

<style scoped>
.typeBtn {
  text-align: left;
  width: 100%;
}
</style>
