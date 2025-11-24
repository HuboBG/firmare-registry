<template>
  <div class="card">
    <h2>Settings</h2>

    <h3>Authentication Mode</h3>
    <p class="small" v-if="!oidcEnabled">
      <strong>API Key Mode (Legacy)</strong><br>
      UI uses hardcoded API keys from environment variables. This is NOT secure for production.<br>
      Consider enabling Keycloak/OIDC for proper authentication.
    </p>
    <p class="small" v-else>
      <strong>Keycloak/OIDC Mode (Secure)</strong><br>
      Authentication is handled via Keycloak. Users must login to access admin features.<br>
      API keys are only used as fallback if no Bearer token is present.
    </p>

    <h3>Keycloak / OIDC Configuration</h3>
    <div v-if="oidcEnabled">
      <p class="small" style="color: green;">✓ OIDC is enabled and configured</p>
      <ul class="small">
        <li><strong>Authority:</strong> {{ oidcAuthority || 'Not set' }}</li>
        <li><strong>Client ID:</strong> {{ oidcClientId || 'Not set' }}</li>
        <li><strong>Redirect URI:</strong> {{ oidcRedirectUri || 'Auto-detected' }}</li>
        <li><strong>Scope:</strong> {{ oidcScope || 'Default' }}</li>
      </ul>
    </div>
    <div v-else>
      <p class="small">To enable OIDC authentication, configure the following:</p>

      <h4>Backend API (.env or environment variables):</h4>
      <ul class="small">
        <li>FW_OIDC_ENABLED=true</li>
        <li>FW_OIDC_ISSUER_URL=https://keycloak.example.com/realms/yourrealm</li>
        <li>FW_OIDC_CLIENT_ID=firmware-admin</li>
        <li>FW_OIDC_ADMIN_ROLE=fw-admin</li>
        <li>FW_OIDC_DEVICE_ROLE=fw-device</li>
        <li>FW_OIDC_AUDIENCE=firmware-registry (optional)</li>
      </ul>

      <h4>Frontend UI (.env or environment variables):</h4>
      <ul class="small">
        <li>VITE_OIDC_ENABLED=true</li>
        <li>VITE_OIDC_AUTHORITY=https://keycloak.example.com/realms/yourrealm</li>
        <li>VITE_OIDC_CLIENT_ID=firmware-admin</li>
        <li>VITE_OIDC_REDIRECT_URI=https://your-app.example.com (optional, auto-detected)</li>
        <li>VITE_OIDC_SCOPE=openid profile email (optional)</li>
      </ul>

      <h4>Keycloak Configuration:</h4>
      <ul class="small">
        <li>Create a client with <strong>Public</strong> access type</li>
        <li>Enable <strong>Standard Flow</strong> (Authorization Code with PKCE)</li>
        <li>Set valid redirect URIs to your application URL</li>
        <li>Create roles: <strong>fw-admin</strong> and <strong>fw-device</strong></li>
        <li>Assign roles to users who need access</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { runtimeConfig } from '../runtime-config';

const oidcEnabled = ref(false);
const oidcAuthority = ref('');
const oidcClientId = ref('');
const oidcRedirectUri = ref('');
const oidcScope = ref('');

onMounted(() => {
  oidcEnabled.value = runtimeConfig.OIDC_ENABLED;
  oidcAuthority.value = runtimeConfig.OIDC_AUTHORITY;
  oidcClientId.value = runtimeConfig.OIDC_CLIENT_ID;
  oidcRedirectUri.value = runtimeConfig.OIDC_REDIRECT_URI;
  oidcScope.value = runtimeConfig.OIDC_SCOPE;
});
</script>
