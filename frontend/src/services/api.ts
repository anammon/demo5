// frontend/src/services/api.ts
import axios from "axios";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:3000";

const api = axios.create({
  baseURL: API_BASE,
  headers: { "Content-Type": "application/json; charset=utf-8" },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers = config.headers || {};
    config.headers["Authorization"] = /^Bearer\s/i.test(token) ? token : `Bearer ${token}`;
  }
  return config;
});

export default api;
