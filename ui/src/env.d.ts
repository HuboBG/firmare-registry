interface ImportMetaEnv {
    readonly VITE_API_BASE_URL?: string;
    readonly VITE_ADMIN_KEY?: string;
    readonly VITE_DEVICE_KEY?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
