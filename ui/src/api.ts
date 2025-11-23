import axios from "axios";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || "",
    timeout: 20000
});

const adminHeaders = () => {
    const k = import.meta.env.VITE_ADMIN_KEY;
    return k ? {"X-Admin-Key": k} : {};
};
const deviceHeaders = () => {
    const k = import.meta.env.VITE_DEVICE_KEY;
    return k ? {"X-Device-Key": k} : {};
};

export interface FirmwareDTO {
    type: string;
    version: string;
    filename: string;
    sizeBytes: number;
    sha256: string;
    createdAt: string;
    downloadUrl?: string;
}

export interface WebhookDTO {
    id: number;
    url: string;
    events: string[];
    enabled: boolean;
}

export const FirmwareAPI = {
    async list(type: string): Promise<FirmwareDTO[]> {
        const r = await api.get(`/api/firmware/${type}`, {headers: deviceHeaders()});
        return r.data as FirmwareDTO[];
    },

    async latest(type: string): Promise<FirmwareDTO> {
        const r = await api.get(`/api/firmware/${type}/latest`, {headers: deviceHeaders()});
        return r.data as FirmwareDTO;
    },

    async upload(type: string, version: string, file: File): Promise<FirmwareDTO> {
        const fd = new FormData();
        fd.append("file", file);
        const r = await api.post(`/api/firmware/${type}/${version}`, fd, {
            headers: {...adminHeaders(), "Content-Type": "multipart/form-data"}
        });
        return r.data as FirmwareDTO;
    },

    async remove(type: string, version: string): Promise<{ deleted: boolean }> {
        const r = await api.delete(`/api/firmware/${type}/${version}`, {headers: adminHeaders()});
        return r.data as { deleted: boolean };
    }
};

export const WebhookAPI = {
    async list(): Promise<WebhookDTO[]> {
        const r = await api.get(`/api/webhooks`, {headers: adminHeaders()});
        return r.data as WebhookDTO[];
    },

    async create(hook: Omit<WebhookDTO, "id">): Promise<{ id: number }> {
        const r = await api.post(`/api/webhooks`, hook, {headers: adminHeaders()});
        return r.data as { id: number };
    },

    async update(id: number, hook: Omit<WebhookDTO, "id">): Promise<{ updated: boolean }> {
        const r = await api.put(`/api/webhooks/${id}`, hook, {headers: adminHeaders()});
        return r.data as { updated: boolean };
    },

    async remove(id: number): Promise<{ deleted: boolean }> {
        const r = await api.delete(`/api/webhooks/${id}`, {headers: adminHeaders()});
        return r.data as { deleted: boolean };
    }
};
