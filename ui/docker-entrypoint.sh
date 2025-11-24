#!/bin/sh
set -e

# Generate runtime config from environment variables
cat > /usr/share/nginx/html/config.js <<EOF
window.ENV = {
  VITE_API_BASE_URL: "${VITE_API_BASE_URL}",
  VITE_ADMIN_KEY: "${VITE_ADMIN_KEY}",
  VITE_DEVICE_KEY: "${VITE_DEVICE_KEY}",
  VITE_OIDC_ENABLED: "${VITE_OIDC_ENABLED}",
  VITE_OIDC_AUTHORITY: "${VITE_OIDC_AUTHORITY}",
  VITE_OIDC_CLIENT_ID: "${VITE_OIDC_CLIENT_ID}",
  VITE_OIDC_REDIRECT_URI: "${VITE_OIDC_REDIRECT_URI}",
  VITE_OIDC_SCOPE: "${VITE_OIDC_SCOPE}"
};
EOF

echo "Runtime configuration generated:"
cat /usr/share/nginx/html/config.js

# Execute the main container command
exec "$@"
