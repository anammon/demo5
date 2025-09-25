// frontend/src/pages/Home.tsx
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { fetchApps, likeApp, getAppLikes } from "../services/app";
import type { AppsPage, AppDTO } from "../services/app";

export default function Home() {
  const navigate = useNavigate();
  const [apps, setApps] = useState<AppDTO[]>([]);
  const [likes, setLikes] = useState<Record<number, number>>({});
  const [q, setQ] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize] = useState(12);
  const [total, setTotal] = useState(0);
  const [inputPage, setInputPage] = useState("");
  const token = localStorage.getItem("token");

  useEffect(() => {
    if (!token) navigate("/login");
  }, [token, navigate]);

  // 拉取当前页数据
  const load = async () => {
    try {
      const res: AppsPage = await fetchApps({ name: q, page, pageSize });
      setApps(res.data || []);
      setTotal(res.total || 0);
    } catch (err) {
      console.error("fetchApps error", err);
    }
  };

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, q]);

  // 获取点赞数（可优化：后端直接返回 likes 字段）
  useEffect(() => {
    apps.forEach((app) => {
      const id = Number(app.ID ?? app.id ?? 0);
      if (!id) return;
      getAppLikes(id).then((r: any) => {
        setLikes((prev) => ({ ...prev, [id]: Number(r.likes || 0) }));
      }).catch(() => {});
    });
  }, [apps]);

  const handleLike = async (app: AppDTO) => {
    const id = Number(app.ID ?? app.id ?? 0);
    if (!id) return;
    try {
      await likeApp(id);
      setLikes((prev) => ({ ...prev, [id]: (prev[id] || 0) + 1 }));
    } catch (err) {
      console.error(err);
    }
  };

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  return (
    <div className="relative w-screen h-screen flex flex-col overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-blue-200 via-blue-100 to-white" />

      {/* 标题 */}
      <div className="relative z-10 text-center py-6">
        <h1 className="text-2xl font-bold text-blue-600">Welcome to AssistApp</h1>
      </div>

      {/* 工具条 */}
      <div className="relative z-10 bg-white/80 backdrop-blur-md shadow px-6 py-4 flex justify-between items-center">
        <input
          placeholder="搜索应用"
          value={q}
          onChange={(e) => { setPage(1); setQ(e.target.value); }}
          className="px-3 py-2 border rounded-md flex-1 max-w-xs"
        />
        <div className="flex items-center gap-3">
          <Link to="/apps/create" className="px-3 py-2 bg-green-600 text-white rounded-md">注册应用</Link>
          <Link to="/apps/update" className="px-3 py-2 bg-indigo-500 text-white rounded-md">更新应用</Link>
          <button onClick={() => { localStorage.removeItem("token"); navigate("/login"); }} className="px-3 py-2 bg-red-500 text-white rounded-md">退出登录</button>
        </div>
      </div>

      {/* 应用网格 */}
      <div className="relative z-10 flex-1 overflow-y-auto p-6">
        <div className="grid gap-6" style={{ gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))" }}>
          {apps.map((app) => {
            const id = Number(app.ID ?? app.id ?? 0);
            return (
              <div key={id} className="bg-white rounded-xl shadow-md p-4 flex flex-col items-center hover:shadow-lg transition">
                {app.Icon || app.icon ? (
                  <img src={app.Icon || app.icon} alt={app.Name || app.name} className="w-20 h-20 mb-3 rounded-lg object-cover" />
                ) : (
                  <div className="w-20 h-20 bg-gray-200 rounded-lg mb-3" />
                )}
                <h3 className="font-semibold text-gray-800 text-center">{app.Name || app.name}</h3>
                <p className="text-sm text-gray-500 line-clamp-2 text-center mb-2">{app.Description || app.description}</p>

                <div className="flex items-center gap-2 mt-auto">
                  <button onClick={() => handleLike(app)} className="text-red-500 hover:scale-110 transition">♥</button>
                  <span className="text-gray-600 text-sm">{likes[id] ?? 0}</span>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* 分页控制 */}
      <div className="relative z-10 bg-white/70 backdrop-blur-sm px-6 py-3 flex justify-center gap-4 items-center">
        <button disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))} className="px-4 py-2 bg-gray-200 rounded-md disabled:opacity-50">上一页</button>
        <span className="text-gray-600">第 {page} / {totalPages} 页</span>
        <button disabled={page >= totalPages} onClick={() => setPage((p) => Math.min(totalPages, p + 1))} className="px-4 py-2 bg-gray-200 rounded-md disabled:opacity-50">下一页</button>

        <input type="number" min={1} max={totalPages} value={inputPage} onChange={(e) => setInputPage(e.target.value)} className="w-16 px-2 py-1 border rounded-md text-center" placeholder="页码" />
        <button onClick={() => { const t = parseInt(inputPage || "0"); if (t >= 1 && t <= totalPages) setPage(t); }} className="px-3 py-1 bg-blue-500 text-white rounded-md">跳转</button>
      </div>
    </div>
  );
}
