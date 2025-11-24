import { UserManager, User, WebStorageStateStore } from 'oidc-client-ts';
import { ref, computed } from 'vue';

export interface AuthConfig {
    enabled: boolean;
    authority: string;
    clientId: string;
    redirectUri: string;
    scope: string;
}

// Global reactive auth state
export const currentUser = ref<User | null>(null);
export const isAuthenticated = computed(() => currentUser.value !== null && !currentUser.value.expired);

let userManager: UserManager | null = null;

export function initAuth(config: AuthConfig) {
    if (!config.enabled) {
        return;
    }

    userManager = new UserManager({
        authority: config.authority,
        client_id: config.clientId,
        redirect_uri: config.redirectUri,
        response_type: 'code',
        scope: config.scope,
        post_logout_redirect_uri: config.redirectUri,
        userStore: new WebStorageStateStore({ store: window.localStorage }),
        automaticSilentRenew: true,
    });

    // Load existing user from storage
    userManager.getUser().then(user => {
        if (user && !user.expired) {
            currentUser.value = user;
        }
    });

    // Handle silent renew success
    userManager.events.addUserLoaded((user) => {
        currentUser.value = user;
    });

    // Handle silent renew error
    userManager.events.addSilentRenewError((error) => {
        console.error('Silent renew error:', error);
    });

    // Handle user signed out
    userManager.events.addUserSignedOut(() => {
        currentUser.value = null;
    });
}

export async function login() {
    if (!userManager) {
        throw new Error('Auth not initialized');
    }
    await userManager.signinRedirect();
}

export async function handleCallback() {
    if (!userManager) {
        throw new Error('Auth not initialized');
    }
    const user = await userManager.signinRedirectCallback();
    currentUser.value = user;
    return user;
}

export async function logout() {
    if (!userManager) {
        throw new Error('Auth not initialized');
    }
    await userManager.signoutRedirect();
}

export function getAccessToken(): string | null {
    return currentUser.value?.access_token || null;
}

export function getUserProfile() {
    return currentUser.value?.profile || null;
}
