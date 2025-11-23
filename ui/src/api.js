import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "",
  timeout: 20000
});

function adminHeaders() {
  const k = import.meta.env.VITE_ADMIN_KEY;
  return k ? { "X-Admin-Key": k } : {};
}
function deviceHeaders() {
  const k = import.meta.env.VITE_DEVICE_KEY;
  return k ? { "X-Device-Key": k } : {};
}

export const FirmwareAPI = {
  async list(type) {
    const r = await api.get(`/api/firmware/${type}`, { headers: deviceHeaders() });
    return r.data;
  },
  async latest(type) {
    const r = await api.get(`/api/firmware/${type}/latest`, { headers: deviceHeaders() });
    return r.data;
  },
  async upload(type, version, file) {
    const fd = new FormData();
    fd.append("file", file);
    const r = await api.post(`/api/firmware/${type}/${version}`, fd, {
      headers: { ...adminHeaders(), "Content-Type": "multipart/form-data" }
    });
    return r.data;
  },
  async remove(type, version) {
    const r = await api.delete(`/api/firmware/${type}/${version}`, { headers: adminHeaders() });
    return r.data;
  }
};

export const WebhookAPI = {
  async list() {
    const r = await api.get(`/api/webhooks`, { headers: adminHeaders() });
    return r.data;
  },
  async create(hook) {
    const r = await api.post(`/api/webhooks`, hook, { headers: adminHeaders() });
    return r.data;
  },
  async update(id, hook) {
    const r = await api.put(`/api/webhooks/${id}`, hook, { headers: adminHeaders() });
    return r.data;
  },
  async remove(id) {
    const r = await api.delete(`/api/webhooks/${id}`, { headers: adminHeaders() });
    return r.data;
  }
};
