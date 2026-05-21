import { state } from "./state.js";

const API_BASE = "/api";

async function request(path, options = {}) {
  const headers = new Headers(options.headers || {});
  headers.set("Content-Type", "application/json");

  if (state.currentUserId) {
    headers.set("X-User-ID", state.currentUserId);
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  const payload = await response.json();
  if (!response.ok || !payload.success) {
    throw new Error(payload.message || "Terjadi kesalahan");
  }

  return payload.data;
}

export const api = {
  getUsers: () => request("/users"),
  getVehicles: () => request("/vehicles"),
  getItems: () => request("/master-items"),
  getReports: (status = "") => request(status ? `/reports?status=${encodeURIComponent(status)}` : "/reports"),
  getReport: (id) => request(`/reports/${id}`),
  createReport: (body) => request("/reports", { method: "POST", body: JSON.stringify(body) }),
  approveReport: (id) => request(`/reports/${id}/approve`, { method: "PATCH" }),
  completeReport: (id, body) => request(`/reports/${id}/complete`, { method: "PATCH", body: JSON.stringify(body) }),
};
