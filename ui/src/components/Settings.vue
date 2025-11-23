<template>
  <div class="card">
    <h2>Settings</h2>
    <p class="small">
      UI uses API keys from env vars. Behind nginx BasicAuth, this is usually enough.
    </p>

    <h3>Keycloak / OIDC (planned)</h3>
    <p class="small">
      API already supports OIDC config fields & middleware stubs.
      When you enable real JWT verification, we can add PKCE login here.
    </p>

    <ul class="small">
      <li>FW_OIDC_ENABLED=true</li>
      <li>FW_OIDC_ISSUER_URL=https://keycloak/realms/yourrealm</li>
      <li>FW_OIDC_CLIENT_ID=firmware-admin</li>
      <li>FW_OIDC_ADMIN_ROLE=fw-admin</li>
      <li>FW_OIDC_DEVICE_ROLE=fw-device</li>
    </ul>
  </div>
</template>
