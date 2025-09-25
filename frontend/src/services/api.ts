// frontend/src/services/api.ts
import axios from "axios";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:3000";

const api = axios.create({
  baseURL: API_BASE,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers = config.headers || {};
    // token 如果已经带 Bearer 前缀，直接用；否则加上
    if (/^Bearer\s/i.test(token)) config.headers["Authorization"] = token;
    else config.headers["Authorization"] = `Bearer ${token}`;
  }
  return config;
});

export default api;
