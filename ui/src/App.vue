<template>
  <div class="app">
    <header class="header">
      <h1>Firmware Registry Admin</h1>
      <nav>
        <button :class="{active:tab==='firmware'}" @click="tab='firmware'">Firmware</button>
        <button :class="{active:tab==='webhooks'}" @click="tab='webhooks'">Webhooks</button>
        <button :class="{active:tab==='settings'}" @click="tab='settings'">Settings</button>
        <button v-if="oidcEnabled && !isAuthenticated" @click="handleLogin" class="login-btn">Login</button>
        <button v-if="oidcEnabled && isAuthenticated" @click="handleLogout" class="logout-btn">
          Logout ({{ userProfile?.name || userProfile?.email || 'User' }})
        </button>
      </nav>
    </header>

    <main v-if="!handlingCallback">
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
    <main v-else>
      <div class="card">
        <p>Completing authentication...</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import {ref, onMounted} from "vue";
import FirmwareTypes from "./components/FirmwareTypes.vue";
import FirmwareVersions from "./components/FirmwareVersions.vue";
import UploadFirmware from "./components/UploadFirmware.vue";
import Webhooks from "./components/Webhooks.vue";
import Settings from "./components/Settings.vue";
import type {FirmwareDTO} from "./api";
import { runtimeConfig } from "./runtime-config";
import { initAuth, login, logout, handleCallback, isAuthenticated, getUserProfile } from "./auth";

const tab = ref<"firmware" | "webhooks" | "settings">("firmware");
const selectedType = ref<string>("");
const oidcEnabled = ref(false);
const handlingCallback = ref(false);
const userProfile = ref<any>(null);

onMounted(async () => {
  // Initialize OIDC if enabled
  oidcEnabled.value = runtimeConfig.OIDC_ENABLED;
  if (oidcEnabled.value) {
    initAuth({
      enabled: true,
      authority: runtimeConfig.OIDC_AUTHORITY,
      clientId: runtimeConfig.OIDC_CLIENT_ID,
      redirectUri: runtimeConfig.OIDC_REDIRECT_URI,
      scope: runtimeConfig.OIDC_SCOPE,
    });

    // Check if this is a callback from Keycloak
    if (window.location.search.includes('code=') || window.location.search.includes('state=')) {
      handlingCallback.value = true;
      try {
        await handleCallback();
        // Clear query params and redirect to home
        window.history.replaceState({}, document.title, window.location.pathname);
      } catch (error) {
        console.error('Authentication callback failed:', error);
      } finally {
        handlingCallback.value = false;
      }
    }

    // Update user profile
    userProfile.value = getUserProfile();
  }
});

// emitted from FirmwareTypes
function selectType(t: string) {
  selectedType.value = t;
}

// currently unused, but handy for future refresh wiring
function onUploaded(_dto: FirmwareDTO) {
  // no-op for now
}

async function handleLogin() {
  await login();
}

async function handleLogout() {
  await logout();
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
