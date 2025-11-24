# Keycloak/OIDC Authentication Setup

This guide explains how to configure Keycloak authentication for the Firmware Registry.

## Overview

The Firmware Registry supports two authentication modes:

1. **API Key Mode (Legacy)**: Hardcoded keys in environment variables - NOT secure for production
2. **Keycloak/OIDC Mode (Recommended)**: Industry-standard OAuth 2.0/OIDC authentication with role-based access control

## Why Use Keycloak?

- **Secure**: No secrets embedded in browser/images
- **User Management**: Centralized user accounts and roles
- **SSO**: Single sign-on across multiple applications
- **Auditing**: Track who accessed what and when
- **MFA**: Multi-factor authentication support

## Architecture

```
User → Browser → Keycloak (login) → Browser gets JWT
Browser → API (with JWT Bearer token) → API verifies JWT → API checks roles
```

## Keycloak Setup

### 1. Create a Realm

1. Login to Keycloak admin console
2. Create a new realm (e.g., `firmware-registry`)
3. Note the realm name - you'll need it for configuration

### 2. Create a Client

1. Go to **Clients** → **Create**
2. Configure:
   - **Client ID**: `firmware-admin` (or your choice)
   - **Client Type**: `OpenID Connect`
   - **Access Type**: `Public` (for browser-based apps)
3. **Settings**:
   - **Valid Redirect URIs**: `https://your-app.example.com/*`
   - **Web Origins**: `https://your-app.example.com`
   - **Standard Flow Enabled**: ON
   - **Direct Access Grants Enabled**: OFF (for security)
4. Save

### 3. Create Roles

1. Go to **Realm Roles** → **Create Role**
2. Create two roles:
   - **fw-admin**: For admin users (upload, delete, webhooks)
   - **fw-device**: For device access (list, download firmware)

### 4. Create Users and Assign Roles

1. Go to **Users** → **Add User**
2. Set username, email, etc.
3. Go to **Credentials** tab → Set password
4. Go to **Role Mapping** tab → Assign roles (`fw-admin` or `fw-device`)

## Backend API Configuration

Set these environment variables in `api/.env`:

```bash
# Enable OIDC
FW_OIDC_ENABLED=true

# Keycloak issuer URL (replace with your realm)
FW_OIDC_ISSUER_URL=https://keycloak.example.com/realms/firmware-registry

# Client ID from Keycloak
FW_OIDC_CLIENT_ID=firmware-admin

# Optional: expected audience claim
FW_OIDC_AUDIENCE=firmware-registry

# Role names (must match Keycloak roles)
FW_OIDC_ADMIN_ROLE=fw-admin
FW_OIDC_DEVICE_ROLE=fw-device

# Keep API keys as fallback (optional)
FW_ADMIN_KEY=your-admin-key
FW_DEVICE_KEY=your-device-key
```

## Frontend UI Configuration

Set these environment variables in `ui/.env`:

```bash
# Enable OIDC
VITE_OIDC_ENABLED=true

# Keycloak authority (same as backend issuer URL)
VITE_OIDC_AUTHORITY=https://keycloak.example.com/realms/firmware-registry

# Client ID (same as backend)
VITE_OIDC_CLIENT_ID=firmware-admin

# Redirect URI (optional, auto-detected if not set)
VITE_OIDC_REDIRECT_URI=https://your-app.example.com

# OIDC scopes (optional, defaults to "openid profile email")
VITE_OIDC_SCOPE=openid profile email

# Keep API keys as fallback (optional)
VITE_ADMIN_KEY=your-admin-key
VITE_DEVICE_KEY=your-device-key
```

## Docker Deployment

When deploying with Docker, add the OIDC variables to your `ui/.env` file:

```bash
VITE_OIDC_ENABLED=true
VITE_OIDC_AUTHORITY=https://keycloak.example.com/realms/firmware-registry
VITE_OIDC_CLIENT_ID=firmware-admin
VITE_OIDC_REDIRECT_URI=https://your-app.example.com
VITE_OIDC_SCOPE=openid profile email
```

The Docker entrypoint script will inject these at runtime into `/config.js`.

## How It Works

### Backend (Go API)

1. API receives request with `Authorization: Bearer <jwt-token>` header
2. API validates JWT signature using Keycloak's JWKS endpoint
3. API extracts roles from JWT claims
4. API checks if user has required role (`fw-admin` or `fw-device`)
5. If valid, request is processed
6. If no Bearer token, falls back to API key check (if configured)

### Frontend (Vue UI)

1. User clicks "Login" button
2. App redirects to Keycloak login page (PKCE flow)
3. User authenticates with username/password
4. Keycloak redirects back with authorization code
5. App exchanges code for JWT token
6. Token is stored in localStorage
7. All API requests include `Authorization: Bearer <token>` header
8. Token is automatically refreshed when expired

## Testing

1. Start backend: `go run ./cmd/firmware-registry`
2. Start frontend: `yarn dev` (development) or deploy with Docker (production)
3. Navigate to the UI
4. Click "Login" - you should be redirected to Keycloak
5. Login with a user that has `fw-admin` role
6. You should be redirected back and see "Logout (username)" button
7. Try uploading firmware - it should work with JWT instead of API keys

## Troubleshooting

### "OIDC enabled but failed to initialize"

- Check that `FW_OIDC_ISSUER_URL` is accessible from the API server
- Verify the issuer URL is correct (should end with `/realms/your-realm`)
- Check Keycloak logs for errors

### "401 Unauthorized" after login

- Verify the user has the correct role assigned in Keycloak
- Check API logs - they should show "User missing required role: fw-admin"
- Verify `FW_OIDC_ADMIN_ROLE` matches the role name in Keycloak

### Login redirect not working

- Check `VITE_OIDC_REDIRECT_URI` matches a valid redirect URI in Keycloak client settings
- Verify the client is set to "Public" access type
- Check browser console for CORS errors

### JWT verification failed

- Verify the JWT is valid at https://jwt.io
- Check that `iss` claim matches `FW_OIDC_ISSUER_URL`
- Check that `aud` claim matches `FW_OIDC_AUDIENCE` (if set)
- Verify Keycloak's JWKS endpoint is accessible: `https://keycloak.example.com/realms/your-realm/protocol/openid-connect/certs`

## Security Best Practices

1. **Use HTTPS**: Never run Keycloak or the app over HTTP in production
2. **Disable API keys**: Once OIDC is working, remove `FW_ADMIN_KEY` and `FW_DEVICE_KEY` from environment variables
3. **Short token lifetime**: Configure Keycloak to issue tokens with short expiration (e.g., 5-15 minutes)
4. **Enable MFA**: Require multi-factor authentication for admin users in Keycloak
5. **Restrict redirect URIs**: Only allow your application's domain in Keycloak client settings
6. **Use roles carefully**: Only assign `fw-admin` to trusted users

## Migration from API Keys

1. Set up Keycloak and configure both backend and frontend
2. Keep API keys configured during migration period (fallback mode)
3. Test OIDC authentication thoroughly
4. Once confirmed working, remove API key environment variables
5. Update documentation for your users to use login instead of API keys
