// frontend/src/services/app.ts
import api from "./api";

export type AppDTO = {
  ID?: number;
  id?: number;
  Name?: string;
  name?: string;
  Description?: string;
  description?: string;
  Icon?: string;
  icon?: string;
  // 其它字段...
};

export type AppsPage = {
  data: AppDTO[];
  total: number;
  page: number;
  pageSize: number;
};

export async function fetchApps(params?: { name?: string; description?: string; page?: number; pageSize?: number }): Promise<AppsPage> {
  const res = await api.get("/app", { params });
  // res.data === { data: [...], total: N, page: P, pageSize: S }
  return res.data as AppsPage;
}

export async function fetchAppById(id: number | string) {
  const res = await api.get(`/app/${id}`);
  return res.data;
}

export async function createApp(payload: Partial<AppDTO>) {
  const res = await api.post("/app", payload);
  return res.data;
}

export async function updateApp(id: number | string, payload: Partial<AppDTO>) {
  const res = await api.put(`/app/${id}`, payload);
  return res.data;
}

export async function likeApp(id: number | string) {
  const res = await api.post(`/app/${id}/like`);
  return res.data;
}

export async function getAppLikes(id: number | string) {
  const res = await api.get(`/app/${id}/likes`);
  return res.data;
}
