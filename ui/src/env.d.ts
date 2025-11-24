interface ImportMetaEnv {
    readonly VITE_API_BASE_URL?: string;
    readonly VITE_ADMIN_KEY?: string;
    readonly VITE_DEVICE_KEY?: string;
    readonly VITE_OIDC_ENABLED?: string;
    readonly VITE_OIDC_AUTHORITY?: string;
    readonly VITE_OIDC_CLIENT_ID?: string;
    readonly VITE_OIDC_REDIRECT_URI?: string;
    readonly VITE_OIDC_SCOPE?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
