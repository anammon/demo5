import axios from "axios";

// 运行时优先使用 VITE_API_BASE（通过构建时注入）。
// 在生产/发布环境中默认使用空字符串（同域），开发可通过 VITE_API_BASE 或本地 dev-server 配置覆盖。
const API_BASE = import.meta.env.VITE_API_BASE ?? "";

const api = axios.create({
  baseURL: API_BASE,
  headers: { "Content-Type": "application/json; charset=utf-8" },
  // 不使用 cookie/session 时关闭 withCredentials，可减少跨域预检复杂性
  withCredentials: false,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers = config.headers || {};
    config.headers["Authorization"] = token.startsWith("Bearer ")
      ? token
      : `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error("API 错误:", error);
    return Promise.reject(error);
  }
);

export default api;
