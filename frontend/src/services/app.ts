// frontend/src/services/app.ts
import api from "./api";

export interface CreateAppDTO {
  name: string;
  description: string;
  icon?: string;
  type?: string;
  tags?: string[];
  status?: string;
  author?: string;
}

export const fetchApps = (q?: { name?: string; description?: string }) =>
  api.get("/app", { params: q }).then((r) => r.data);

export const fetchAppById = (id: number | string) =>
  api.get(`/app/${id}`).then((r) => r.data);

export const createApp = (data: CreateAppDTO) =>
  api.post("/app", data).then((r) => r.data);

export const likeApp = (id: number | string) =>
  api.post(`/app/${id}/like`).then((r) => r.data);

export const getAppLikes = (id: number | string) =>
  api.get(`/app/${id}/likes`).then((r) => r.data);
